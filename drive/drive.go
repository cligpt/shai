package drive

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/cligpt/shai/config"
	rpc "github.com/cligpt/shai/drive/rpc"
)

const (
	upTiemout = "10s"
)

type Drive interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, string, string) (string, string, error)
}

type Config struct {
	Logger hclog.Logger
	Config config.Config
}

type drive struct {
	cfg    *Config
	client rpc.RpcProtoClient
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

func (d *drive) Init(ctx context.Context) error {
	if err := d.initConn(ctx); err != nil {
		return errors.Wrap(err, "failed to init conn")
	}

	return nil
}

func (d *drive) Deinit(ctx context.Context) error {
	_ = d.deinitConn(ctx)

	return nil
}

func (d *drive) Run(ctx context.Context, role, content string) (r, c string, e error) {
	ret, err := d.sendChat(ctx, role, content)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to chat")
	}

	return ret.GetMessage().GetRole(), ret.GetMessage().GetContent(), nil
}

func (d *drive) initConn(_ context.Context) error {
	var err error

	host := d.cfg.Config.Spec.Drive.Host
	port := d.cfg.Config.Spec.Drive.Port

	d.conn, err = grpc.NewClient(host+":"+strconv.Itoa(port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)))
	if err != nil {
		return errors.Wrap(err, "failed to dial")
	}

	d.client = rpc.NewRpcProtoClient(d.conn)

	return nil
}

func (d *drive) deinitConn(_ context.Context) error {
	return d.conn.Close()
}

func (d *drive) sendChat(ctx context.Context, role, content string) (*rpc.ChatReply, error) {
	ctx, cancel := context.WithTimeout(ctx, d.setTimeout(upTiemout))
	defer cancel()

	// TBD: FIXME
	messages := []*rpc.ChatMessage{
		{
			Role:    role,
			Content: content,
		},
	}

	reply, err := d.client.SendChat(ctx, &rpc.ChatRequest{
		Model:    "llama3",
		Messages: messages,
		Format:   "json",
		Options: &rpc.ChatOption{
			Temperature: 0,
		},
		Stream:    false,
		KeepAlive: "5m",
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to send")
	}

	return reply, nil
}

func (d *drive) setTimeout(timeout string) time.Duration {
	duration, _ := time.ParseDuration(timeout)

	return duration
}
