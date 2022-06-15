package pbgo

import (
	"golang.org/x/net/context"
	"io"
)

type HelloStudentServiceServerImp struct {
}

func (p *HelloStudentServiceServerImp) Hello(ctx context.Context, args *Student) (*Student, error) {
	reply := &Student{Name: args.GetName()}
	return reply, nil
}
func (p *HelloStudentServiceServerImp) Channel(stream HelloStudentService_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &Student{Name: args.GetName()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}
