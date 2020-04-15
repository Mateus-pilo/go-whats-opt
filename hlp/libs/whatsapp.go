package libs

import (
	"encoding/base64"
	"encoding/gob"
	"errors"
	"log"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"
	"net/http"
	"encoding/json"
	"path/filepath"
	"bytes"

	
	
	whatsapp "github.com/Rhymen/go-whatsapp"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/Mateus-pilo/go-whats-opt/hlp"
)

var wac = make(map[string]*whatsapp.Conn)



type waHandler struct{
	c *whatsapp.Conn
	jid string
	created uint64
}


type msgResponse struct {
	whatsapp.TextMessage
	Jid string `json:"jid"`
}

type msgResponseImage struct {
	whatsapp.ImageMessage
	Type  string
	Jid string `json:"jid"`
}

type msgResponseDocument struct {
	whatsapp.DocumentMessage
	Type  string
	Jid string `json:"jid"`
}


type msgResponseVideo struct {
	whatsapp.VideoMessage
	Type string
	Jid string `json:"jid"`
}
type msgResponseAudio struct {
	whatsapp.AudioMessage
	Type string
	Jid string `json:"jid"`
}

type responseContacts struct {
	Contacts []whatsapp.Contact `json:"contacts"`
	Jid string `json:"jid_company"`
}




func (h *waHandler) HandleError(err error) {

	if e, ok := err.(*whatsapp.ErrConnectionFailed); ok {
		log.Printf("Connection failed, underlying error: %v", e.Err)
		log.Println("Waiting 30sec...")
		<-time.After(30 * time.Second)
		log.Println("Reconnecting...")
		err := h.c.Restore()
		if err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
	} else {
		log.Printf("error occoured: %v\n", err)
	}
}


func (h *waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	if message.Info.FromMe == false && message.Info.Timestamp >= h.created  {
		
		responseMessage := msgResponse{TextMessage: message}
		responseMessage.Jid = h.jid
		
		jsonStr, _ := json.Marshal(responseMessage)
		urlPost := hlp.Config.GetString("SERVER_API_NODE")
		//fmt.Println(urlPost);
		req, _ := http.NewRequest("POST", urlPost, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
	}
}

func (h *waHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	
	
	if message.Info.FromMe == false && message.Info.Timestamp >= h.created {
		data, err := message.Download()
		if err != nil {
			if err != whatsapp.ErrMediaDownloadFailedWith410 {
				return
			}
			data, err = message.Download()
			if err != nil {
				return
			
			}
		}
		filename := fmt.Sprintf("%v/%v.%v", "/var/whats/image", message.Info.Id, strings.Split(message.Type, "/")[1])
		file, err := os.Create(filename)
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			fmt.Printf("[!] %v\n", err)
			return
		}

		
		body := new(bytes.Buffer)

		writer := multipart.NewWriter(body)

		_, err = writer.CreateFormFile("file", filepath.Base(filename))

    if err != nil {
        log.Fatal(err)
		}

		responseMessage := msgResponseImage{ImageMessage: message, Type: "image", Jid: h.jid}
		

		var inInterface map[string]string
		jsonStr, _ := json.Marshal(responseMessage)
		json.Unmarshal(jsonStr, &inInterface)
		
    // iterate through inrecs
    for field, val := range inInterface {
			if field != "Info" {
				_ = writer.WriteField(field, val)
			}
		}
		
		var inInterfaceInfo map[string]string
		jsonStr, _ = json.Marshal(message.Info)
    json.Unmarshal(jsonStr, &inInterfaceInfo)

		for field, val := range inInterfaceInfo {
			_ = writer.WriteField(field, val)
		}
		

		writer.Close()
		

		urlPost := hlp.Config.GetString("SERVER_API_NODE")
		req, _ := http.NewRequest("POST", urlPost, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()

	}
	
}

func (h *waHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	
	if message.Info.FromMe == false  && message.Info.Timestamp >= h.created {
	
		data, err := message.Download()
		if err != nil {
			if err != whatsapp.ErrMediaDownloadFailedWith410 {
				return
			}
			data, err = message.Download()
			if err != nil {
				return
			
			}
		}
		filename := fmt.Sprintf("%v/%v.%v", "/var/whats/document", message.Info.Id, strings.Split(message.Type, "/")[1])
		file, err := os.Create(filename)
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			fmt.Printf("[!] %v\n", err)
			return
		}

		
		body := new(bytes.Buffer)

		writer := multipart.NewWriter(body)

		_, err = writer.CreateFormFile("file", filepath.Base(filename))

    if err != nil {
        log.Fatal(err)
		}

		responseMessage := msgResponseDocument{DocumentMessage: message, Type: "file", Jid: h.jid}
		

		var inInterface map[string]string
		jsonStr, _ := json.Marshal(responseMessage)
		json.Unmarshal(jsonStr, &inInterface)
		
    // iterate through inrecs
    for field, val := range inInterface {
			if field != "Info" {
				_ = writer.WriteField(field, val)
			}
		}
		
		var inInterfaceInfo map[string]string
		jsonStr, _ = json.Marshal(message.Info)
    json.Unmarshal(jsonStr, &inInterfaceInfo)

    // iterate through inrecs
    for field, val := range inInterfaceInfo {
			_ = writer.WriteField(field, val)
		}
		
		writer.Close()
		

		urlPost := hlp.Config.GetString("SERVER_API_NODE")
		req, _ := http.NewRequest("POST", urlPost, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
	}
}

func (h *waHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	
	if message.Info.FromMe == false  && message.Info.Timestamp >= h.created {
	
		data, err := message.Download()
		if err != nil {
			if err != whatsapp.ErrMediaDownloadFailedWith410 {
				return
			}
			data, err = message.Download()
			if err != nil {
				return
			
			}
		}
		filename := fmt.Sprintf("%v/%v.%v", "/var/whats/video", message.Info.Id, strings.Split(message.Type, "/")[1])
		file, err := os.Create(filename)
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			fmt.Printf("[!] %v\n", err)
			return
		}

		
		body := new(bytes.Buffer)

		writer := multipart.NewWriter(body)

		_, err = writer.CreateFormFile("file", filepath.Base(filename))

    if err != nil {
        log.Fatal(err)
		}

		responseMessage := msgResponseVideo{VideoMessage: message, Type: "video", Jid: h.jid}
		

		var inInterface map[string]string
		jsonStr, _ := json.Marshal(responseMessage)
		json.Unmarshal(jsonStr, &inInterface)
		
    // iterate through inrecs
    for field, val := range inInterface {
			if field != "Info" {
				_ = writer.WriteField(field, val)
			}
		}
		
		var inInterfaceInfo map[string]string
		jsonStr, _ = json.Marshal(message.Info)
    json.Unmarshal(jsonStr, &inInterfaceInfo)

    // iterate through inrecs
    for field, val := range inInterfaceInfo {
			_ = writer.WriteField(field, val)
		}

		writer.Close()
		

		urlPost := hlp.Config.GetString("SERVER_API_NODE")
		req, _ := http.NewRequest("POST", urlPost, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
	}
}

func (h *waHandler) HandleAudioMessage(message whatsapp.AudioMessage){	
	
	if message.Info.FromMe == false  && message.Info.Timestamp >= h.created {
	
		data, err := message.Download()
		if err != nil {
			if err != whatsapp.ErrMediaDownloadFailedWith410 {
				return
			}
			data, err = message.Download()
			if err != nil {
				return
			
			}
		}
		filename := fmt.Sprintf("%v/%v.%v", "/var/whats/audio", message.Info.Id, strings.Split(message.Type, "/")[1])
		file, err := os.Create(filename)
		defer file.Close()

		_, err = file.Write(data)
		if err != nil {
			fmt.Printf("[!] %v\n", err)
			return
		}

		
		body := new(bytes.Buffer)

		writer := multipart.NewWriter(body)

		_, err = writer.CreateFormFile("file", filepath.Base(filename))

    if err != nil {
        log.Fatal(err)
		}

		responseMessage := msgResponseAudio{AudioMessage: message, Type: "audio", Jid: h.jid}
		

		var inInterface map[string]string
		jsonStr, _ := json.Marshal(responseMessage)
		json.Unmarshal(jsonStr, &inInterface)
		
    // iterate through inrecs
    for field, val := range inInterface {
			if field != "Info" {
				_ = writer.WriteField(field, val)
			}
		}
		
		var inInterfaceInfo map[string]string
		jsonStr, _ = json.Marshal(message.Info)
    json.Unmarshal(jsonStr, &inInterfaceInfo)

    // iterate through inrecs
    for field, val := range inInterfaceInfo {
			_ = writer.WriteField(field, val)
		}
		writer.Close()

		urlPost := hlp.Config.GetString("SERVER_API_NODE")
		req, _ := http.NewRequest("POST", urlPost, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		client := &http.Client{}
		resp, _ := client.Do(req)
		defer resp.Body.Close()
	}
}


func (h *waHandler) HandleContactList(Contacts[] whatsapp.Contact){
	responseContact := responseContacts{Contacts: Contacts}
	responseContact.Jid = h.jid
	
	jsonStr, _ := json.Marshal(responseContact)
	urlPost := hlp.Config.GetString("SERVER_API_NODE_CONTACTS")
	req, _ := http.NewRequest("POST", urlPost, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()	
}

func (h *waHandler) HandleChatList(contacts[] whatsapp.Chat){
	
	//fmt.Println(contacts);
}


func WASyncVersion(conn *whatsapp.Conn) (string, error) {
	conn.SetClientVersion(0, 4,  2081)
	versionClient := conn.GetClientVersion()
	
	return fmt.Sprintf("whatsapp version %v.%v.%v", versionClient[0], versionClient[1], versionClient[2]), nil
}

func WATestPing(conn *whatsapp.Conn) error {
	ok, err := conn.AdminTest()
	if !ok {
		if err != nil {
			return err
		} else {
			return errors.New("something when wrong while trying to ping, please check phone connectivity")
		}
	}

	return nil
}

func WAGenerateQR(timeout int, chanqr chan string, qrstr chan<- string) {
	select {
	case tmp := <-chanqr:
		png, _ := qrcode.Encode(tmp, qrcode.Medium, 256)
		qrstr <- base64.StdEncoding.EncodeToString(png)
	}
}

func WASessionInit(jid string, timeout int) error {
	if wac[jid] == nil {
		conn, err := whatsapp.NewConn(time.Duration(timeout) * time.Second)
		if err != nil {
			return err
		}
		conn.SetClientVersion(0, 4,  2081)
		conn.SetClientName("Hiper Chat", "Go Whats")

		info, err := WASyncVersion(conn)
		if err != nil {
			return err
		}
		hlp.LogPrintln(hlp.LogLevelInfo, "whatsapp", info)

		
		wac[jid] = conn

	}

	return nil
}

func WASessionLoad(file string) (whatsapp.Session, error) {
	session := whatsapp.Session{}

	buffer, err := os.Open(file)
	if err != nil {
		return session, err
	}
	defer buffer.Close()

	err = gob.NewDecoder(buffer).Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

func WASessionSave(file string, session whatsapp.Session) error {
	buffer, err := os.Create(file)
	if err != nil {
		return err
	}
	defer buffer.Close()

	err = gob.NewEncoder(buffer).Encode(session)
	if err != nil {
		return err
	}

	return nil
}

func WASessionExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	return true
}

func WASessionConnect(jid string, timeout int, file string, qrstr chan<- string, errmsg chan<- error) {
	chanqr := make(chan string)

	session, err := WASessionLoad(file)
	if err != nil {
		go func() {
			WAGenerateQR(timeout, chanqr, qrstr)
		}()
		
		err = WASessionLogin(jid, timeout, file, chanqr)
		if err != nil {
			errmsg <- err
			return
		}
		return
	}

	err = WASessionRestore(jid, timeout, file, session)
	if err != nil {
		go func() {
			WAGenerateQR(timeout, chanqr, qrstr)
		}()

		err = WASessionLogin(jid, timeout, file, chanqr)
		if err != nil {
			errmsg <- err
			return
		}
	}

	err = WATestPing(wac[jid])
	if err != nil {
		errmsg <- err
		return
	}

	errmsg <- errors.New("")
	return
}

func WASessionLogin(jid string, timeout int, file string, qrstr chan<- string) error {
	if wac[jid] != nil {
		if WASessionExist(file) {
			err := os.Remove(file)
			if err != nil {
				return err
			}
		}

		delete(wac, jid)
	}

	err := WASessionInit(jid, timeout)
	if err != nil {
		return err
	}

	session, err := wac[jid].Login(qrstr)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "already logged in":
			return nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return errors.New("connection is invalid")
		default:
			delete(wac, jid)
			return err
		}
	}

	err = WASessionSave(file, session)
	if err != nil {
		return err
	}

	wac[jid].AddHandler(&waHandler{wac[jid], jid, uint64(time.Now().Unix())})

	return nil
}

func WASessionRestore(jid string, timeout int, file string, sess whatsapp.Session) error {
	if wac[jid] != nil {
		if WASessionExist(file) {
			err := os.Remove(file)
			if err != nil {
				return err
			}
		}

		delete(wac, jid)
	}

	err := WASessionInit(jid, timeout)
	if err != nil {
		return err
	}

	session, err := wac[jid].RestoreWithSession(sess)
	if err != nil {
		switch strings.ToLower(err.Error()) {
		case "already logged in":
			return nil
		case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
			delete(wac, jid)
			return errors.New("connection is invalid")
		default:
			delete(wac, jid)
			return err
		}
	}

	err = WASessionSave(file, session)
	if err != nil {
		return err
	}

	wac[jid].AddHandler(&waHandler{wac[jid], jid, uint64(time.Now().Unix())})
	
	return nil
}

func WASessionLogout(jid string, file string) error {
	if wac[jid] != nil {
		err := wac[jid].Logout()
		if err != nil {
			return err
		}

		if WASessionExist(file) {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}

		delete(wac, jid)
	} else {
		return errors.New("connection is invalid")
	}

	return nil
}

func WAMessageText(jid string, jidDest string, msgText string, msgQuotedID string, msgQuoted string, msgDelay int) (string, error) {
	var id string
	if wac[jid] != nil {
		jidPrefix := "@s.whatsapp.net"
		if len(strings.SplitN(jidDest, "-", 2)) == 2 {
			jidPrefix = "@g.us"
		}

		content := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: jidDest + jidPrefix,
			},
			Text: msgText,
		}

		<-time.After(time.Duration(msgDelay) * time.Second)
		id, err := wac[jid].Send(content)
		if err != nil {
			switch strings.ToLower(err.Error()) {
			case "sending message timed out":
				return id, nil
			case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
				delete(wac, jid)
				return "", errors.New("connection is invalid")
			default:
				return "", err
			}
		}
	} else {
		return "", errors.New("connection is invalid")
	}

	return id, nil
}

func WAMessageImage(jid string, jidDest string, msgImageStream multipart.File, msgImageType string, msgCaption string, msgDelay int) (string, error) {
	var id string

	if wac[jid] != nil {
		jidPrefix := "@s.whatsapp.net"
		if len(strings.SplitN(jidDest, "-", 2)) == 2 {
			jidPrefix = "@g.us"
		}

		content := whatsapp.ImageMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: jidDest + jidPrefix,
			},
			Content: msgImageStream,
			Type:    msgImageType,
			Caption: msgCaption,
		}

		<-time.After(time.Duration(msgDelay) * time.Second)

		id, err := wac[jid].Send(content)
		if err != nil {
			switch strings.ToLower(err.Error()) {
			case "sending message timed out":
				return id, nil
			case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
				delete(wac, jid)
				return "", errors.New("connection is invalid")
			default:
				return "", err
			}
		}
	} else {
		return "", errors.New("connection is invalid")
	}
	return id, nil
}

func WAMessageDocument(jid string, jidDest string, msgDocumentStream multipart.File, msgDocumentType string, msgDocumentInfo string, msgDelay int) (string, error) {
	var id string

	if wac[jid] != nil {
		jidPrefix := "@s.whatsapp.net"
		if len(strings.SplitN(jidDest, "-", 2)) == 2 {
			jidPrefix = "@g.us"
		}

		content := whatsapp.DocumentMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: jidDest + jidPrefix,
			},
			Content:  msgDocumentStream,
			Type:     msgDocumentType,
			FileName: msgDocumentInfo,
			Title:    msgDocumentInfo,
		}

		<-time.After(time.Duration(msgDelay) * time.Second)

		id, err := wac[jid].Send(content)
		if err != nil {
			switch strings.ToLower(err.Error()) {
			case "sending message timed out":
				return id, nil
			case "could not send proto: failed to write message: error writing to websocket: websocket: close sent":
				delete(wac, jid)
				return "", errors.New("connection is invalid")
			default:
				return "", err
			}
		}
	} else {
		return "", errors.New("connection is invalid")
	}

	return id, nil
}
