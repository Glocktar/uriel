package main

import (
	"strings"
	"strconv"
	"fmt"
	"time"
	"io/ioutil"
	"github.com/bwmarrin/discordgo"
	"net/http"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	conf = Configuration{}
	warning bool
)

type Configuration struct {
	dSession 	*discordgo.Session
	dChannel	string
	dBotID 		string
	dToken		string
}

func main() {
	var err error

	//Configuration Loader
	filebuff, err := ioutil.ReadFile("uriel_config.ini")
	if err != nil {
		fmt.Println("No configuration file found.")
		panic(err)
	}

	fileconts := strings.Split(string(filebuff),"\r\n")
	if len(fileconts) < 2 {
		fmt.Println("No configuration file found.")
		panic(err)
	}
	conf.dToken = fileconts[0]
	conf.dChannel = fileconts[1]
	// ----------------------------

	conf.dSession, err = discordgo.New("Bot " + conf.dToken)
	check(err)
	
	u, err := conf.dSession.User("@me")
	check(err)
	
	conf.dBotID = u.ID
	conf.dSession.AddHandler(commands)

	err = conf.dSession.Open()
	check(err)

	fmt.Println("Uriel logged in.")

	warning = false

	SetTimer(5 * time.Minute, Updater)
	
	<-make(chan struct{})
	return
}

func commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == conf.dBotID {
		return
	}
	
	var cmd, params, params2, params3 string
	fmt.Sscanf(m.Content,"%s %s %s %s",&cmd,&params,&params2,&params3)

	if cmd == "!cmds" || cmd == "!cmd" || cmd == "!comandos" || cmd == "!commands" {
		_, _ = s.ChannelMessageSend(conf.dChannel,"Comandos Uriel:\n!time : Muestra la hora del Uriel\n!hora : Muestra la hora de España\n!noticias : Muestra las noticias mas recientes de la página de Metin2.es\n!eventos : Muestra los eventos fijos del juego")
	}

	if cmd == "!time" {
		_, _ = s.ChannelMessageSend(conf.dChannel,"Hora actual Venezuela : " + time.Now().Format(time.RFC850))
	}

	if cmd == "!hora" {
		resp, err := http.Get("https://24timezones.com/es_husohorario/madrid_hora_actual.php")
		check(err)
	
		root, err := html.Parse(resp.Body)
		check(err)
	
		matcher := func(n *html.Node) bool {
			if n.DataAtom == atom.Span {
				return scrape.Attr(n, "id") == "currentTime"
			}
			return false
		}
	
		messagem, _ := scrape.Find(root,matcher)
		_, _ = s.ChannelMessageSend(conf.dChannel,"Hora actual CEST: " + scrape.Text(messagem))
	}

	if cmd == "!noticias" || cmd == "!news" {
		resp, err := http.Get("https://board.es.metin2.gameforge.com/")
		check(err)

		root, err := html.Parse(resp.Body)
		check(err)

		matcher := func(n *html.Node) bool {
			if n.DataAtom == atom.Header || n.DataAtom == atom.Div {
				return scrape.Attr(n,"class") == "messageHeader" || scrape.Attr(n,"class") == "messageText"
			}
			return false
		}

		articles := scrape.FindAll(root, matcher)
		_, _ = s.ChannelMessageSend(conf.dChannel,"**Noticias Metin2.es**")
		_, _ = s.ChannelMessageSend(conf.dChannel,"**" + scrape.Text(articles[0]) + "**" + "\n" + scrape.Text(articles[1]))
		_, _ = s.ChannelMessageSend(conf.dChannel,"**" + scrape.Text(articles[2]) + "**" + "\n" + scrape.Text(articles[3]))
		_, _ = s.ChannelMessageSend(conf.dChannel,"**" + scrape.Text(articles[4]) + "**" + "\n" + scrape.Text(articles[5]))
		_, _ = s.ChannelMessageSend(conf.dChannel,"**" + scrape.Text(articles[6]) + "**" + "\n" + scrape.Text(articles[7]))
		_, _ = s.ChannelMessageSend(conf.dChannel,"**" + scrape.Text(articles[8]) + "**" + "\n" + scrape.Text(articles[9]))
	}


	if cmd == "!eventos" || cmd == "!event" || cmd == "!events" || cmd == "!evento" {
		_, _ = s.ChannelMessageSend(conf.dChannel,"**Eventos Metin2.es**")
		_, _ = s.ChannelMessageSend(conf.dChannel,"**Lunes:**\n\t22:00-00:00 Drop de supermonturas\n\n**Martes:**\n\t16:00-18:00 Drop de supermonturas\n\n**Miércoles:**\n\t23:00-01:00 Drop de supermonturas\n\n**Jueves:**\n\t13:00-15:00 Drop de supermonturas\n\n**Viernes:**\n\t18:00-20:00 Drop de supermonturas\n\n**Sábado:**\n\t13:00-15:00 Drop de supermonturas\n\t18:00-00:00 (El segundo sábado de cada mes) Drop de cajas luz luna\n\tTodo el día (Último fin de semana de cada mes) Festival de la cosecha\n\n**Domingo:**\n\t13:00-15:00 Drop de supermonturas\n\t22:00-02:00 Drop de Alubias verdes del dragón\n\tTodo el día (Último fin de semana de cada mes) Festival de la cosecha\n\n**Cualquier día de la semana:**\n\tCompetición OX (Mínimo una vez por semana)\n\tDrop de telas delicadas 4 h. (Sin día ni hora predefinido, sale el mensaje en el juego cuando comienza)\n\tDrop de Cor Draconis 4 h. (Sin día ni hora predefinido, sale el mensaje en el juego cuando comienza).")
	}

}

func SetTimer(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

func Updater() {
	t := cestTime()

	if strings.Contains(t.dayna,"Lunes") {
		if t.hour >= 10 && t.apm == "PM" {
			if warning == false {
				_, _ = conf.dSession.ChannelMessageSend(conf.dChannel,"@everyone \nEl drop de supermonturas ha empezado y terminará a las 12 PM. Hora actual CEST: " + strconv.Itoa(t.hour) + ":" + strconv.Itoa(t.min))
				warning = true
			}
		} else {
			warning = false
		}
	}

	if strings.Contains(t.dayna,"Martes") {
		if t.hour >= 4 && t.hour < 6 && t.apm == "PM" {
			if warning == false {
				_, _ = conf.dSession.ChannelMessageSend(conf.dChannel,"@everyone \nEl drop de supermonturas ha empezado y terminará a las 6 PM. Hora actual CEST: " + strconv.Itoa(t.hour) + ":" + strconv.Itoa(t.min))
				warning = true				
			}
		} else {
			warning = false
		}
	}

	if strings.Contains(t.dayna,"Mi") {
		if t.hour >= 23 && t.apm == "PM" {
			if warning == false {
				_, _ = conf.dSession.ChannelMessageSend(conf.dChannel,"@everyone \nEl drop de supermonturas ha empezado y terminará a las 1 AM. Hora actual CEST: " + strconv.Itoa(t.hour) + ":" + strconv.Itoa(t.min))
				warning = true				
			}
		} else {
			warning = false
		}
	}

	if strings.Contains(t.dayna,"Jueves") {
		if t.hour >= 1 && t.hour < 3 && t.apm == "PM" {
			if warning == false {
				_, _ = conf.dSession.ChannelMessageSend(conf.dChannel,"@everyone \nEl drop de supermonturas ha empezado y terminará a las 3 PM. Hora actual CEST: " + strconv.Itoa(t.hour) + ":" + strconv.Itoa(t.min))
				warning = true				
			}
		} else {
			warning = false
		}
	}	

	if strings.Contains(t.dayna,"Viernes") {
		if t.hour >= 6 && t.hour < 8 && t.apm == "PM" {
			if warning == false {
				_, _ = conf.dSession.ChannelMessageSend(conf.dChannel,"@everyone \nEl drop de supermonturas ha empezado y terminará a las 8 PM. Hora actual CEST: " + strconv.Itoa(t.hour) + ":" + strconv.Itoa(t.min))
				warning = true				
			}
		} else {
			warning = false
		}
	}

	if strings.Contains(t.dayna,"bado") {
		if t.hour >= 1 && t.hour < 3 && t.apm == "PM" {
			if warning == false {
				_, _ = conf.dSession.ChannelMessageSend(conf.dChannel,"@everyone \nEl Festival de la Cosecha ha comenzado.\nEl drop de supermonturas ha empezado y terminará a las 3 PM. Hora actual CEST: " + strconv.Itoa(t.hour) + ":" + strconv.Itoa(t.min))
				warning = true
			}
		} else {
			warning = false
		}
	}	

	if strings.Contains(t.dayna,"Domingo") {
		if t.hour >= 1 && t.hour < 3 && t.apm == "PM" {
			if warning == false {
				_, _ = conf.dSession.ChannelMessageSend(conf.dChannel,"@everyone \nEl drop de supermonturas ha empezado y terminará a las 3 PM. Hora actual CEST: " + strconv.Itoa(t.hour) + ":" + strconv.Itoa(t.min))
				warning = true				
			}
		} else {
			warning = false
		}
	}	

}

type Taimu struct {
	apm		string
	hour	int
	min		int
	sec		int
	daynu	int
	dayna	string
	monthna	string
}

func cestTime() Taimu {
	t := Taimu{}

	resp, err := http.Get("https://24timezones.com/es_husohorario/madrid_hora_actual.php")
	check(err)

	root, err := html.Parse(resp.Body)
	check(err)

	matcher := func(n *html.Node) bool {
		if n.DataAtom == atom.Span {
			return scrape.Attr(n, "id") == "currentTime"
		}
		return false
	}

	var hourstr string

	timestr, _ := scrape.Find(root,matcher)

	fmt.Sscanf(strings.Replace(scrape.Text(timestr), ",","",-1),"%s %s %s %d %s", &hourstr, &t.apm, &t.dayna, &t.daynu, &t.monthna) //04:41:39 PM, Domingo 10, septiembre 2017

	hourarray := strings.Split(hourstr,":")

	t.hour, _ = strconv.Atoi(hourarray[0])
	t.min, _ = strconv.Atoi(hourarray[1])
	t.sec, _ = strconv.Atoi(hourarray[2])

	return t
}

func serverTime() Taimu {
	timenow := time.Now().Format(time.RFC850)
	taimu := Taimu{}

	var hourstr, datestr string

	fmt.Sscanf(timenow,"%s %s %5s",&taimu.dayna, &datestr, &hourstr)//Sunday, 10-Sep-17 09:45:28 -04
	
	hourminsec := strings.Split(hourstr,":")
	daymonthyear := strings.Split(datestr,"-")

	taimu.hour, _ = strconv.Atoi(hourminsec[0])
	taimu.min, _ = strconv.Atoi(hourminsec[1])
	taimu.sec, _ = strconv.Atoi(hourminsec[2])
	taimu.daynu, _ = strconv.Atoi(daymonthyear[0])
	taimu.monthna = daymonthyear[1]

	return taimu
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
