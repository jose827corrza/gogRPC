package main

import (
	"context"
	"github.com/jose827corrza/gogRPC/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

/*
Este codigo hace lo que realiza Postman, osea.. consume el back tanto Unary como Stream desde el server
*/
func main() {
	cc, err := grpc.Dial("localhost:5020", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer cc.Close()

	c := testpb.NewTestServiceClient(cc)
	//DoUnary(c)
	//DoCLientStream(c)
	//DoServerStreamin(c)
	DoBidirectionalStreamin(c)
}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{Id: "t1"}

	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GetTest: %v", err)
	}
	log.Printf("response from server: %v", res)
}

func DoCLientStream(c testpb.TestServiceClient) {
	questions := []*testpb.Question{
		{
			Id:       "q1t1",
			TestId:   "t1",
			Question: "Especialidad de Go",
			Answer:   "Backend",
		},
		{
			Id:       "q1t1",
			TestId:   "t1",
			Question: "Lenguaje favorito de Jose",
			Answer:   "Javascript",
		},
	}

	stream, err := c.SetQuestion(context.Background())
	if err != nil {
		log.Fatalf("Error while calling SetQuestion: %v", err)
	}
	for _, question := range questions {
		log.Printf("sending question: %v", question.Id)
		stream.Send(question)
		time.Sleep(2 * time.Second)
	}
	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}
	log.Printf("response from server: %v", msg)
}

func DoServerStreamin(c testpb.TestServiceClient) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "t1",
	}
	stream, err := c.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GetStudentPerTest: %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from erver: %v", msg)
	}
}

func DoBidirectionalStreamin(c testpb.TestServiceClient) {
	answer := testpb.TakeTestRequest{
		Answer: "42",
	}
	numberOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := c.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("Error while calling TakeTest: %v", err)
	}
	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			stream.Send(&answer)
			time.Sleep(2 * time.Second)
		}
	}()
	go func() {
		for {
			stream, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while reading streamin: %v", err)
				break
			}
			log.Printf("response from server: %v", stream)
		}
		close(waitChannel)
	}()
	<-waitChannel //Con esto se bloquea el programa
}
