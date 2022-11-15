package main

//Программа читает настройки из clisnt.json
//Устанавливает TCP соединение с сервером из настроек
//Держит соединение, пока жив сервер
//Постоянно читает файл command.txt, если он есть, выполняет содержимое и переименовывает файл в _command.txt

import (
	"fmt"
	"net"
	"encoding/binary"
	"context"
	"time"	
	"io"
	"io/ioutil"
	"strings"
	"os"
)

const (
	PREF_PACKET_START byte = 0xFF	//Обязательный префикс начала пакета
	PREF_PACKET_LAST byte = 0x0A	//Префикс последнего пакета, или единственного
	PREF_PACKET_CONT byte = 0x0B	//Префикс пакета, закоторым последуют еще пакеты, тут будет только часть сообщения
	POSTF_0 byte = 0x0A		//Постфикс
	POSTF_1 byte = 0x0D 		//Постфикс
	PREF_LEN uint32 = 10
	POSTF_LEN uint32 = 2
	
)

var packet_id uint32

func main() {

	/*if len(os.Args)<2 {
		panic("Не найден аргумент с командой")
	}*/

	//чтение файла настроек
	conf := AppConfig{}	
	if err := conf.ReadConf("client.json"); err != nil {
		panic(fmt.Sprintf("failed conf.ReadConf: %v",err))
	}
	
	//Соединение, таймаут - 1 минута
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	
	//соединение с данными из настройки
	//conf.Server = IP:PORT
	conn, err := d.DialContext(ctx, "tcp", conf.Server)
	if err != nil {
		panic(fmt.Sprintf("Failed to dial: %v", err))
	}
	defer conn.Close()

	//сообщение для отправки из параметров командной строки
	//msg := `{"func":"Test.insert", "argv": {"f1":99, "f2":"Еще немного текста!","f4":true}}`
	//msg := os.Args[1] //из параметров командной строки
	//sendToServer(conn, msg)

	//проверка файла на существование, если есть - исполнение команды
	go sendToServerFromFile(conn)

	//вечный цикл ожидания ответа пока
	readFromServer(conn)
}

func sendToServerFromFile(conn net.Conn){
	for {
		file, err := ioutil.ReadFile("command.txt")
		if err == nil {
			
			lines := strings.Split(string(file),"\n")
			for _,s := range lines {				
				if s != "" {
					//fmt.Println(s)
					sendToServer(conn, s)
				}
			}
			
			os.Rename("command.txt", "_command.txt")
		}	
		
		time.Sleep(2 * time.Second)
	}
}

//Отправка сообщения на сервер
//Сообщение делится на куски по 65535 байт максимум
//Типы сообщений (второй байт от начала) могут быть: PREF_PACKET_LAST/PREF_PACKET_CONT
func sendToServer(conn net.Conn, msg string) {

	packet_len := uint32(len(msg))	
	
	bf := make([]byte, PREF_LEN + packet_len + POSTF_LEN)
	bf[0] = PREF_PACKET_START
	bf[1] = PREF_PACKET_LAST
	binary.LittleEndian.PutUint32(bf[2:6], packet_len)
	binary.LittleEndian.PutUint32(bf[6:10], packet_id)		
	copy(bf[PREF_LEN : PREF_LEN+packet_len], msg)
	bf[PREF_LEN+packet_len] = POSTF_0
	bf[PREF_LEN+packet_len+1] = POSTF_1
	
	//bf := append([]byte{PREF_PACKET_START}, packet_type, packet_len_b, msg[data_offset:data_offset+packet_len], []byte{POSTF_0}, []byte{POSTF_0}...)
	_, err := conn.Write(bf)
	if err != nil {
		panic(fmt.Sprintf("Failed conn.Write(): %v", err))
		return err		
	}
	packet_id++

	/*
	//Вся длина сообщения
	whole_len := len(msg)
	fmt.Printf("Общая длина сообщения: %d, %s\n",whole_len, msg)
	
	//Количество пакетов для отправки
	packet_cnt := whole_len / int(MAX_DATA_LEN) + 1
	fmt.Printf("Сообщение разделили на %d пакетов\n",packet_cnt)
	
	//смещение от начала данных
	var data_offset uint16
	
	//цикл по пакетам
	for packet_n := 0; packet_n < packet_cnt; packet_n++{
		//Это тип пакета: PREF_PACKET_LAST/PREF_PACKET_CONT, т.е. или заевршающий пакет или с продолжением
		var packet_type byte
		
		//Это длина текущего пакета
		var packet_len uint16
		
		if packet_n+1 == packet_cnt {
			//Последний пакет
			packet_type = PREF_PACKET_LAST
			packet_len = uint16(whole_len - int(data_offset))
		}else{
			//пакет с продолжением
			packet_type = PREF_PACKET_CONT
			packet_len = MAX_DATA_LEN		
		}
		
		//Буфер под пакет: Префикс(2 bytes) + длина данных(2 bytes) + данные JSON (=data length), Постфикс(2 bytes)
		fmt.Printf("Буфер сообщения длина байт: %d\n", 2 + 2 + packet_len +2)
		bf := make([]byte, 2 + 2 + packet_len +2)
		
		//первые два байта 
		bf[0] = PREF_PACKET_START //Всегда
		bf[1] = packet_type //0x0A Если без разбивки, или 0x0B если это часть
		
		//Длина пакета 2-3 байты по правилу Little Endian (младший байт вперед)
		binary.LittleEndian.PutUint16(bf[2:4], packet_len)
		
		//Сами данные (или часть), начиная с 4 байта по 4+packet_len
		//Либо часть сообщения либо все
		copy(bf[4:4+packet_len], msg[data_offset:data_offset+packet_len])
		
		//постфикс
		bf[4+packet_len] = POSTF_0
		bf[4+packet_len+1] = POSTF_1
		
		//отправка
		fmt.Printf("Отправляем пакет %d\n",packet_n)
		//fmt.Println(bf) //отладка - содаржимое буфера
		_, err := conn.Write(bf)
		if err != nil {
			panic(fmt.Sprintf("Failed conn.Write(): %v", err))
		}
		data_offset += packet_len
	}
	*/	
}

func readFromServer(conn net.Conn) {
	//Буфер заголовка
	head_b := make([]byte, 2 + 2)
	
	for {				
		_, err := conn.Read(head_b)			
		switch err {		
		case nil:
			//prefix check
			if head_b[0] != PREF_PACKET_START || (head_b[1] != PREF_PACKET_LAST && head_b[1] != PREF_PACKET_CONT) {
				//wrong structute
				fmt.Println("TCPServer.HandleConnection() wrong packet structure")
				continue			
			}

			//Packet structure:
			//PREFIX(2 bytes) + data length(2 bytes) + JSON data (=data length), POSTF(2 bytes)
		
			packet_len := binary.LittleEndian.Uint16(head_b[2:4])
			if packet_len > MAX_DATA_LEN {
				fmt.Println("read max message length exceeded")
				continue
			}
			
			payload := make([]byte, packet_len+2) //Data + postfix
			_, err := conn.Read(payload)			
			switch err {
			case nil:
				//got message
				
				//postfix check
				if payload[packet_len] != POSTF_0 || payload[packet_len+1] != POSTF_1 {
					fmt.Println("read wrong packet structure")
					continue
				}
				
				//все кроме постфикса
				fmt.Println("Ответ сервера:", string(payload[:packet_len]))
							
			case io.EOF:
				fmt.Println("conn.Read closed")
				
			default:
				fmt.Printf("failed conn.Read: %v\n", err)
			}		
					
		case io.EOF:
			fmt.Println("conn.Read closed")
			return
			
		default:
			fmt.Printf("failed conn.Read: %v\n", err)
			return
		}		
		
	}
}

