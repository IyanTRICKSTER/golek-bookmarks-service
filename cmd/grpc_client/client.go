package grpc_client

import (
	"context"
	"errors"
	"fmt"
	"golek_bookmark_service/pkg/contracts"
	"golek_bookmark_service/pkg/models"
	ps "golek_bookmark_service/pkg/models/proto_schema"
	"google.golang.org/grpc"
	"log"
)

type GRPCServiceClient struct {
	HOST   string
	PORT   string
	Client ps.PostServiceClient
}

func (c *GRPCServiceClient) Fetch(ctx context.Context, postIDs []string) ([]models.Post, error) {

	//cID := ps.CoursesID{CoursesID: []string{"6300988647b1637e7974b3d9", "6300988647b1637e7974b3d6"}}
	courses, err := c.Client.Fetch(ctx, &ps.PostIDs{Id: postIDs})
	if err != nil {
		log.Println("gRPC Client: PostService: Fetch Error >>", err)
		return nil, err
	}

	posts := make([]models.Post, 0)
	for _, c := range courses.List {
		posts = append(posts, models.Post{ID: models.GenerateObjectIDFromHex(c.Id), Name: c.Name})
	}

	return posts, nil
}

func (c *GRPCServiceClient) Dial() (ps.PostServiceClient, error) {

	host := c.HOST + ":" + c.PORT

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("could not connect to %v %v", host, err))
	}

	log.Println("GRPC Connected to", host)

	c.Client = ps.NewPostServiceClient(conn)

	return c.Client, nil
}

func New(config contracts.Config) contracts.GRPCPostService {
	return &GRPCServiceClient{HOST: config.GetAppConfig()["RPC_TARGET_HOST"], PORT: config.GetAppConfig()["RPC_TARGET_PORT"]}
}

//Procedural Test
//func serviceCourse() ps.CoursesServiceClient {
//
//	host := "172.24.0.3:6060"
//
//	conn, err := grpc.Dial(host, grpc.WithInsecure())
//	if err != nil {
//		log.Fatal("could not connect to", host, err)
//	}
//
//	return ps.NewCoursesServiceClient(conn)
//}
//
//func main() {
//
//	coursesServices := serviceCourse()
//
//	cID := ps.CoursesID{CoursesID: []string{"6300988647b1637e7974b3d9", "6300988647b1637e7974b3d6"}}
//
//	courses, err := coursesServices.Fetch(context.TODO(), &cID)
//	if err != nil {
//		return
//	}
//
//	for _, c := range courses.Fetch {
//		log.Println(c)
//	}
//
//}
