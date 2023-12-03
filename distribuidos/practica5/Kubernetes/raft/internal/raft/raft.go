// Escribir vuestro código de funcionalidad Raft en este fichero
//

package raft

//
// API
// ===
// Este es el API que vuestra implementación debe exportar
//
// nodoRaft = NuevoNodo(...)
//   Crear un nuevo servidor del grupo de elección.
//
// nodoRaft.Para()
//   Solicitar la parado de un servidor
//
// nodo.ObtenerEstado() (yo, mandato, esLider)
//   Solicitar a un nodo de elección por "yo", su mandato en curso,
//   y si piensa que es el msmo el lider
//
// nodoRaft.SometerOperacion(operacion interface()) (indice, mandato, esLider)

// type AplicaOperacion

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	//"crypto/rand"
	"math/rand"
	"sync"
	"time"

	//"net/rpc"

	"raft/internal/comun/rpctimeout"
)

const (
	// Constante para fijar valor entero no inicializado
	IntNOINICIALIZADO = -1

	//  false deshabilita por completo los logs de depuracion
	// Aseguraros de poner kEnableDebugLogs a false antes de la entrega
	kEnableDebugLogs = true

	// Poner a true para logear a stdout en lugar de a fichero
	kLogToStdout = false

	// Cambiar esto para salida de logs en un directorio diferente
	kLogOutputDir = "./logs_raft/"
)

type TipoOperacion struct {
	Operacion string // La operaciones posibles son "leer" y "escribir"
	Clave     string
	Valor     string // en el caso de la lectura Valor = ""
}

// A medida que el nodo Raft conoce las operaciones de las  entradas de registro
// comprometidas, envía un AplicaOperacion, con cada una de ellas, al canal
// "canalAplicar" (funcion NuevoNodo) de la maquina de estados
type AplicaOperacion struct {
	Indice    int // en la entrada de registro
	Operacion TipoOperacion
}

//Strcut que representa la informacion de una entrada
type Entrada struct {
	Index int
	Term  int
	Op    TipoOperacion
}

// Tipo de dato Go que representa un solo nodo (réplica) de raft
//
type NodoRaft struct {
	Mux sync.Mutex // Mutex para proteger acceso a estado compartido

	// Host:Port de todos los nodos (réplicas) Raft, en mismo orden
	Nodos   []rpctimeout.HostPort
	Yo      int // indice de este nodos en campo array "nodos"
	IdLider int
	estado  string //El nodo podra ser seguidor, candidato o lider
	votos   int    //votos recibidos <---------------------------- Quitar

	// Utilización opcional de este logger para depuración
	// Cada nodo Raft tiene su propio registro de trazas (logs)
	Logger *log.Logger

	// Vuestros datos aqui.
	currentTerm int       //Ultimo mandato
	voteFor     int       //Candidato que recibio voto en el mandato actual
	log         []Entrada //log entries

	commitIndex int // Indice del último valor comprometido
	lastApplied int // Indice del último valor aplicado a la maquina de estados

	NextIndex  []int // Para cada nodo la siguiente entrada que tiene que enviar el lider
	MatchIndex []int // Para cada nodo el mayor indice raplicado conocido

	canalLider       chan bool //Indicativo que es lider
	canalSeguidor    chan bool //Indicativo que es seguidor
	latido           chan bool //Latidos que da el lider para indicar que esta vivo
	AplicarOperacion chan AplicaOperacion
}

func tiempoEsperaAleatorio() time.Duration {
	// Seed para generar números realmente aleatorios
	rand.Seed(time.Now().UnixNano())

	// Generar un número aleatorio entre 100 y 200
	return time.Duration(rand.Intn(101)+100) * time.Millisecond
}

func maquinaEstadosNodo(nr *NodoRaft) {
	for {

		if nr.lastApplied < nr.commitIndex {
			fmt.Printf("LastApplied es %d y commitIndex es %d\n", nr.lastApplied, nr.commitIndex)
			nr.AplicarOperacion <- AplicaOperacion{nr.lastApplied,
				nr.log[nr.lastApplied].Op}
			nr.lastApplied++
		}

		if nr.estado == "seguidor" {
			select {
			case <-nr.latido:

			case <-time.After(tiempoEsperaAleatorio()):
				nr.IdLider = -1
				nr.estado = "candidato"
			}
		} else if nr.estado == "candidato" {
			pedirVotacion(nr)
			select {
			case <-nr.canalLider:
				nr.estado = "lider"
			case <-nr.canalSeguidor:
				nr.estado = "seguidor"
			case <-nr.latido:
				nr.estado = "seguidor"
			case <-time.After(tiempoEsperaAleatorio()):

			}
		} else if nr.estado == "lider" {
			nr.IdLider = nr.Yo
			nr.mantenerCommitIndex()
			enviarLatidos(nr)
			select {
			case <-nr.canalSeguidor:
				nr.estado = "seguidor"
			case <-time.After(time.Duration(50 * time.Millisecond)):

			}
		}
	}
}

//------------------------------------------------------------------------
// Mantenemos up-to-date
func (nr *NodoRaft) mantenerCommitIndex() {
	exito := 0
	for i := 0; i < len(nr.Nodos); i++ {
		if nr.MatchIndex[i] > nr.commitIndex {
			exito++
		}
		if exito > len(nr.Nodos)/2 {
			nr.commitIndex++
			fmt.Printf("Commit Index ahora es %d\n", nr.commitIndex)
			fmt.Printf("Commit Index ahora es %d\n", nr.commitIndex)
			break
		}
	}
}

// Creacion de un nuevo nodo de eleccion
//
// Tabla de <Direccion IP:puerto> de cada nodo incluido a si mismo.
//
// <Direccion IP:puerto> de este nodo esta en nodos[yo]
//
// Todos los arrays nodos[] de los nodos tienen el mismo orden

// canalAplicar es un canal donde, en la practica 5, se recogerán las
// operaciones a aplicar a la máquina de estados. Se puede asumir que
// este canal se consumira de forma continúa.
//
// NuevoNodo() debe devolver resultado rápido, por lo que se deberían
// poner en marcha Gorutinas para trabajos de larga duracion
func NuevoNodo(nodos []rpctimeout.HostPort, yo int,
	canalAplicarOperacion chan AplicaOperacion) *NodoRaft {
	time.Sleep(5000 * time.Millisecond) //Espera a que se pongan los servidores en ready
	nr := &NodoRaft{}
	nr.Nodos = nodos
	nr.Yo = yo
	nr.IdLider = -1
	nr.estado = "seguidor"
	nr.votos = 0
	nr.currentTerm = 0
	nr.voteFor = IntNOINICIALIZADO
	nr.commitIndex = 0
	nr.lastApplied = 0
	nr.canalLider = make(chan bool)
	nr.canalSeguidor = make(chan bool)
	nr.latido = make(chan bool)
	nr.AplicarOperacion = canalAplicarOperacion
	nr.log = []Entrada{}

	if kEnableDebugLogs {
		nombreNodo := nodos[yo].Host() + "_" + nodos[yo].Port()
		logPrefix := fmt.Sprintf("%s", nombreNodo)

		fmt.Println("LogPrefix: ", logPrefix)

		if kLogToStdout {
			nr.Logger = log.New(os.Stdout, nombreNodo+" -->> ",
				log.Lmicroseconds|log.Lshortfile)
		} else {
			err := os.MkdirAll(kLogOutputDir, os.ModePerm)
			if err != nil {
				panic(err.Error())
			}
			logOutputFile, err := os.OpenFile(fmt.Sprintf("%s/%s.txt",
				kLogOutputDir, logPrefix), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				panic(err.Error())
			}
			nr.Logger = log.New(logOutputFile,
				logPrefix+" -> ", log.Lmicroseconds|log.Lshortfile)
		}
		nr.Logger.Println("logger initialized")
	} else {
		nr.Logger = log.New(ioutil.Discard, "", 0)
	}

	fmt.Printf("Inicializarmos el log con longitud: %d\n", len(nr.log))

	go maquinaEstadosNodo(nr)

	return nr
}

// Metodo Para() utilizado cuando no se necesita mas al nodo
//
// Quizas interesante desactivar la salida de depuracion
// de este nodo
//
func (nr *NodoRaft) para() {
	go func() { time.Sleep(5 * time.Millisecond); os.Exit(0) }()
}

// Devuelve "yo", mandato en curso y si este nodo cree ser lider
//
// Primer valor devuelto es el indice de este  nodo Raft el el conjunto de nodos
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) obtenerEstado() (int, int, bool, int) {
	var yo int = nr.Yo
	var mandato int
	var esLider bool
	var idLider int = nr.IdLider

	mandato = nr.currentTerm
	esLider = idLider == yo

	return yo, mandato, esLider, idLider
}

// El servicio que utilice Raft (base de datos clave/valor, por ejemplo)
// Quiere buscar un acuerdo de posicion en registro para siguiente operacion
// solicitada por cliente.

// Si el nodo no es el lider, devolver falso
// Sino, comenzar la operacion de consenso sobre la operacion y devolver en
// cuanto se consiga
//
// No hay garantia que esta operacion consiga comprometerse en una entrada de
// de registro, dado que el lider puede fallar y la entrada ser reemplazada
// en el futuro.
// Primer valor devuelto es el indice del registro donde se va a colocar
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) someterOperacion(operacion TipoOperacion) (int, int,
	bool, int, string) {
	//nr.Mux.Lock()

	indice := len(nr.log) + 1
	mandato := nr.currentTerm
	EsLider := nr.Yo == nr.IdLider
	idLider := -1
	valorADevolver := ""

	fmt.Println("Ha entrado en someterOperación")

	if EsLider {

		entrada := Entrada{indice, mandato, operacion}

		nr.log = append(nr.log, entrada)

		nr.verLog()

		nr.NextIndex[nr.Yo] = len(nr.log) + 1

		idLider = nr.Yo
	}

	//nr.Mux.Unlock()

	return indice, mandato, EsLider, idLider, valorADevolver
}

//------------------------------------------------------------------------
// Para el test 4
func (nr *NodoRaft) obtenerEstadoLog() (int, int, int, TipoOperacion) {
	return nr.Yo, nr.log[len(nr.log)-1].Index,
		nr.log[len(nr.log)-1].Term,
		nr.log[len(nr.log)-1].Op
}

// -----------------------------------------------------------------------
// LLAMADAS RPC al API
//
// Si no tenemos argumentos o respuesta estructura vacia (tamaño cero)
type Vacio struct{}

func (nr *NodoRaft) ParaNodo(args Vacio, reply *Vacio) error {
	defer nr.para()
	return nil
}

type EstadoParcial struct {
	Mandato int
	EsLider bool
	IdLider int
}

type EstadoRemoto struct {
	IdNodo int
	EstadoParcial
}

func (nr *NodoRaft) ObtenerEstadoNodo(args Vacio, reply *EstadoRemoto) error {
	reply.IdNodo, reply.Mandato, reply.EsLider, reply.IdLider = nr.obtenerEstado()
	return nil
}

type ResultadoRemoto struct {
	ValorADevolver string
	IndiceRegistro int
	EstadoParcial
}

func (nr *NodoRaft) SometerOperacionRaft(operacion TipoOperacion,
	reply *ResultadoRemoto) error {
	reply.IndiceRegistro, reply.Mandato, reply.EsLider, reply.IdLider, reply.ValorADevolver = nr.someterOperacion(operacion)
	return nil
}

//-------------------------------------------------------------------------
// Para el test 4
type EstadoLog struct {
	IdNodo int
	Entrada
}

func (nr *NodoRaft) ObtenerEstadoLogRaft(args Vacio, reply *EstadoLog) error {
	reply.IdNodo, reply.Index, reply.Term, reply.Op = nr.obtenerEstadoLog()
	return nil
}

// -----------------------------------------------------------------------
// LLAMADAS RPC protocolo RAFT
//
// Structura de ejemplo de argumentos de RPC PedirVoto.
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
//
type ArgsPeticionVoto struct {
	Term         int // Mandato del candidato
	CandidateId  int // Candidato solicitando voto
	LastLogIndex int
	LastLogTerm  int
}

// Structura de ejemplo de respuesta de RPC PedirVoto,
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
//
//
type RespuestaPeticionVoto struct {
	Term        int  // Mandato actual, los candidatos lo actualizan
	VoteGranted bool // true significa que el candidato ha recibido mi voto
}

// Metodo para RPC PedirVoto
//
func (nr *NodoRaft) PedirVoto(peticion *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) error {
	nr.Mux.Lock()
	
	lastLogIndex := len(nr.log)
	lastLogTerm := 0
	if lastLogIndex > 0 {
		lastLogTerm = nr.log[lastLogIndex-1].Term
	}
	fmt.Printf("LastLogTerm %d , LastLogIndex %d\n", peticion.LastLogTerm, peticion.LastLogIndex)
	
	nr.Mux.Unlock()

	if peticion.Term > nr.currentTerm {
		nr.canalSeguidor <- true
		nr.currentTerm = peticion.Term
		nr.voteFor = -1
	}

	if peticion.Term == nr.currentTerm && 
	(nr.voteFor == IntNOINICIALIZADO || nr.voteFor == peticion.CandidateId) &&
	(peticion.LastLogTerm > lastLogTerm || (peticion.LastLogTerm == lastLogTerm && peticion.LastLogIndex >= lastLogIndex)){
		reply.VoteGranted = true
		nr.voteFor = peticion.CandidateId
	} else {
		reply.VoteGranted = false
	}

	reply.Term = nr.currentTerm

	return nil
}

type ArgAppendEntries struct {
	Term         int // Mandato del lider
	LeaderId     int // Para que el seguidor pueda redirigir a los clientes
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []Entrada
	LeaderCommit int
}

type Results struct {
	Term    int
	Success bool
}

// Metodo de tratamiento de llamadas RPC AppendEntries
func (nr *NodoRaft) AppendEntries(args *ArgAppendEntries,
	results *Results) error {

	nr.Mux.Lock()

	nr.latido <- true

	if args.Term < nr.currentTerm {
		results.Term = nr.currentTerm
		results.Success = false
		nr.Mux.Unlock()
		return nil
	}

	if len(args.Entries) > 0 {
		fmt.Printf("PrevLogIndex: %d y Len: %d\n", args.PrevLogIndex, len(nr.log))
		results.Success = false

		if args.PrevLogIndex == 0 || (args.PrevLogIndex <= len(nr.log) && args.PrevLogTerm == nr.log[args.PrevLogIndex-1].Term) {
			results.Success = true

			InsertIndex := args.PrevLogIndex
			newEntries := 0

			for {
				if InsertIndex >= len(nr.log) || newEntries >= len(args.Entries) {
					break
				}
				if nr.log[InsertIndex].Term != args.Entries[newEntries].Term {
					break
				}
				InsertIndex++
				newEntries++
			}

			if newEntries < len(args.Entries) {
				nr.log = append(nr.log[:InsertIndex], args.Entries[newEntries:]...)
			}

			if args.LeaderCommit > nr.commitIndex {
				nr.commitIndex = min(args.LeaderCommit, len(nr.log))
			}

			nr.verLog()
		}
	}

	results.Term = nr.currentTerm

	nr.Mux.Unlock()

	return nil
}

func min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

// ----- Metodos/Funciones a utilizar como clientes
//
//

// Ejemplo de código enviarPeticionVoto
//
// nodo int -- indice del servidor destino en nr.nodos[]
//
// args *RequestVoteArgs -- argumentos par la llamada RPC
//
// reply *RequestVoteReply -- respuesta RPC
//
// Los tipos de argumentos y respuesta pasados a CallTimeout deben ser
// los mismos que los argumentos declarados en el metodo de tratamiento
// de la llamada (incluido si son punteros
//
// Si en la llamada RPC, la respuesta llega en un intervalo de tiempo,
// la funcion devuelve true, sino devuelve false
//
// la llamada RPC deberia tener un timout adecuado.
//
// Un resultado falso podria ser causado por una replica caida,
// un servidor vivo que no es alcanzable (por problemas de red ?),
// una petición perdida, o una respuesta perdida
//
// Para problemas con funcionamiento de RPC, comprobar que la primera letra
// del nombre  todo los campos de la estructura (y sus subestructuras)
// pasadas como parametros en las llamadas RPC es una mayuscula,
// Y que la estructura de recuperacion de resultado sea un puntero a estructura
// y no la estructura misma.
//
func (nr *NodoRaft) enviarPeticionVoto(nodo int, args *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) bool {

	fallo := nr.Nodos[nodo].CallTimeout("NodoRaft.PedirVoto", args,
		reply, 50*time.Millisecond)

	if fallo == nil {
		//En el caso que se pida voto a un mandato superior
		if reply.Term > nr.currentTerm {
			nr.currentTerm = reply.Term
			nr.voteFor = -1
			nr.canalSeguidor <- true

		} else if reply.VoteGranted { //Se recibe voto
			
			nr.votos++
			if nr.votos > (len(nr.Nodos) / 2) {
				//Tiene mayoria por lo que se proclama lider
				nr.NextIndex = make([]int, 3)
				nr.MatchIndex = make([]int, 3)

				for i := 0; i < len(nr.Nodos); i++ {
					// Inicializamos nextIndex y matchIndex
					fmt.Println("Soy lider")

					nr.Mux.Lock()
					nr.NextIndex[i] = len(nr.log) + 1
					nr.MatchIndex[i] = 0
					nr.Mux.Unlock()

					fmt.Printf("Para el nodo: %d NextIndex: %d y len: %d\n", i, nr.NextIndex[i], len(nr.log))

				}

				nr.canalLider <- true
			}
			
		}
		return true

	} else {
		return false
	}
}

func pedirVotacion(nr *NodoRaft) {
	var respuesta RespuestaPeticionVoto
	lastLogIndex := 0
	lastLogTerm := 0
	nr.voteFor = nr.Yo
	nr.votos = 1
	nr.currentTerm++

	nr.Mux.Lock()
	if len(nr.log) > 0 {
		lastLogIndex = len(nr.log)
		lastLogTerm = nr.log[lastLogIndex-1].Term
	}
	nr.Mux.Unlock()

	for i := 0; i < len(nr.Nodos); i++ {
		if nr.Yo != i {
			go nr.enviarPeticionVoto(i, &ArgsPeticionVoto{nr.currentTerm, nr.Yo,
				lastLogIndex, lastLogTerm},
				&respuesta)
		}
	}
}

func (nr *NodoRaft) enviarLatido(nodo int, args ArgAppendEntries) bool {

	var reply Results

	fallo := nr.Nodos[nodo].CallTimeout("NodoRaft.AppendEntries", args,
		&reply, 50*time.Millisecond)

	if fallo == nil {
		nr.Mux.Lock()
		//En el caso que se pida voto a un mandato superior
		if reply.Term > nr.currentTerm {
			nr.currentTerm = reply.Term
			nr.IdLider = -1
			nr.voteFor = -1
			nr.canalSeguidor <- true
		}

		if len(args.Entries) > 0 {
			if reply.Success {
				lenght := len(nr.log)
				nr.NextIndex[nodo] = lenght + 1
				nr.MatchIndex[nodo] = lenght
			} else {
				fmt.Println("")
				fmt.Printf("Para el nodo %d\n", nodo)
				fmt.Printf("Se ha restado NextIndex, antes: %d\n", nr.NextIndex[nodo])
				if nr.NextIndex[nodo] > 1 {
					nr.NextIndex[nodo]--
				}
				fmt.Printf("Nuevo NextIndex: %d\n", nr.NextIndex[nodo])
				fmt.Println("")
			}
		}
		nr.Mux.Unlock()
		return true

	} else {
		fmt.Printf("La llamada ha fallado para %d\n", nodo)
		return false
	}
}

func enviarLatidos(nr *NodoRaft) {

	fmt.Println("Soy lider")

	for i := 0; i < len(nr.Nodos); i++ {
		if nr.Yo != i {
			nr.Mux.Lock()
			prevLogIndex := 0
			prevLogTerm := 0

			if lenght := len(nr.log); lenght > 0 {
				prevLogIndex = lenght
				prevLogTerm = nr.log[prevLogIndex-1].Term
			}

			entrada := []Entrada{}

			fmt.Printf("Longitud es: %d y NextIndex es: %d para nodo %d\n", len(nr.log), nr.NextIndex[i], i)

			if len(nr.log) >= nr.NextIndex[i] {

				fmt.Printf("Se envia una entrada para %d\n", i)

				if nr.NextIndex[i] < 2 {
					prevLogIndex = 0
					prevLogTerm = 0
					entrada = nr.log[0:]
				} else {
					prevLogIndex = nr.NextIndex[i] - 1
					prevLogTerm = nr.log[prevLogIndex-1].Term
					entrada = nr.log[nr.NextIndex[i]-1:]
				}

				//nr.Logger.Printf("Para el nodo %d, NextIndex: %d, Len: %d and PrevLogIndex: %d", i, nr.NextIndex[i], len(nr.log), PrevLogIndex)
			}
			nr.Mux.Unlock()
			nr.enviarLatido(i,
				ArgAppendEntries{nr.currentTerm,
					nr.Yo, prevLogIndex,
					prevLogTerm, entrada,
					nr.commitIndex})
		}
	}

}

func (nr *NodoRaft) verLog() {
	fmt.Printf("El Logger con longitud %d contiene\n", len(nr.log))
	if len(nr.log) > 0 {
		for i := 0; i < len(nr.log); i++ {
			fmt.Printf("%d %s\n", nr.log[i].Index, nr.log[i].Op.Operacion)
		}
	} else {
		fmt.Println("El logger esta vacio")
	}

	fmt.Println("")
}
