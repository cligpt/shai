package drive

import (
	"context"
	"math"
	"strconv"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cligpt/shai/config"
	rpc "github.com/cligpt/shai/drive/rpc"
)

type Drive interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context) error
}

type Config struct {
	Logger hclog.Logger
	Config config.Config
}

type drive struct {
	cfg    *Config
	client rpc.AiProtoClient
	conn   *grpc.ClientConn
}

func New(_ context.Context, cfg *Config) Drive {
	return &drive{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (d *drive) Init(_ context.Context) error {
	var err error

	host := d.cfg.Config.Host
	port := d.cfg.Config.Port

	d.conn, err = grpc.NewClient(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	d.client = rpc.NewAiProtoClient(d.conn)

	return nil
}

func (d *drive) Deinit(_ context.Context) error {
	return d.conn.Close()
}

func (d *drive) Run(_ context.Context) error {
	return nil
}
