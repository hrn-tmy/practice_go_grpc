package main

import (
	"fmt"
	"go-grpc/pb"
	"log"
	"os"

	// 旧式："github.com/golang/protobuf/jsonpb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func main() {
	employee := &pb.Employee{
		Id:         1,
		Name:       "Suzuki",
		Email:      "test@example.com",
		Occupation: pb.Occupation_ENGINEER,
		Phone:      []string{"080-1234-5678", "090-1234-5678"},
		Project:    map[string]*pb.Company_Project{"ProjextX": &pb.Company_Project{}},
		Profile: &pb.Employee_Text{
			Text: "My Name Is Suzuki",
		},
		Birthday: &pb.Date{
			Year:  2000,
			Month: 1,
			Day:   1,
		},
	}
	// バイナリ形式にシリアライズする
	binData, err := proto.Marshal(employee)
	if err != nil {
		log.Fatalln("Can't serialize", err)
	}
	// ioutil.WriteFile()は旧式の書き方なので、osパッケージのWriteFileを使うようにする
	if err := os.WriteFile("test.bin", binData, 0666); err != nil {
		log.Fatalln("Can't Write", err)
	}
	// ioutil.ReadFile()は旧式の書き方なので、osパッケージのReadFileを使うようにする
	in, err := os.ReadFile("test.bin")
	if err != nil {
		log.Fatalln("Can't read file", err)
	}

	// バイナリ形式のファイルをデシリアライズする
	readEmployee := &pb.Employee{}
	err = proto.Unmarshal(in, readEmployee)
	if err != nil {
		log.Fatalln("Can't read file", err)
	}
	fmt.Println(readEmployee)

	// JSON形式にシリアライズ
	out, err := protojson.Marshal(employee)
	if err != nil {
		log.Fatalln("Can't marshal to json", err)
	}
	fmt.Println(string(out))

	// JSON形式からデシリアライズ
	jsomReadEmployee := &pb.Employee{}
	if err := protojson.Unmarshal(out,jsomReadEmployee); err != nil {
		log.Println("Can't unmarshal from json", err)
	}
	fmt.Println(jsomReadEmployee)

	// jsonpbは旧式の書き方なので、protojsonを使うようにする
	// m := jsonpb.Marshaler{}
	// out, err := m.MarshalToString(employee)
	// if err != nil {
	// 	log.Fatalln("Can't marshal to json", err)
	// }
	// fmt.Println(out)

	// jsomReadEmployee := &pb.Employee{}
	// if err := jsonpb.UnmarshalString(out,jsomReadEmployee); err != nil {
	// 	log.Println("Can't unmarshal from json", err)
	// }
	// fmt.Println(jsomReadEmployee)
}
