//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"social-app/internal/config"
	"social-app/internal/connector"
	"github.com/google/wire"
	"social-app/internal/domains/media"
	"social-app/internal/domains/profile"
	"social-app/internal/domains/notification"
	"social-app/internal/domains/like"
	"social-app/internal/domains/post"
	"social-app/internal/domains/comment"
	"social-app/internal/domains/chat"
	"social-app/internal/domains/auth"
	"social-app/pkg/ws"
	"social-app/pkg/server"
	"social-app/internal/routes"
	"social-app/pkg/notifier"
	"social-app/internal/domains/llm"
)

var (
	CommonWiring = wire.NewSet(
		config.NewConfig,
		wire.FieldsOf(
			new(*config.Config),
			"DB",
			"Auth",
			"Redis",
			"SMS",
			"Mailtrap",
			"LLM",
		),
		connector.NewDBConn,
	)
	NotifierWiring = wire.NewSet(
		notifier.NewSMS,
		notifier.NewMailTrap,

		wire.Bind(new(notifier.Mailerx), new(*notifier.MailTrap)),
		wire.Bind(new(notifier.SMSNotifier), new(*notifier.SMS)),
	)
	WSWiring = wire.NewSet(
		connector.NewRedisConnector,
		ws.NewConnector,
	)
	RepositoryWiring = wire.NewSet(
		auth.NewRepository,
		chat.NewRepository,
		comment.NewRepository,
		post.NewRepository,
		like.NewRepository,
		media.NewRepository,
		notification.NewRepository,
		profile.NewRepository,
	)
	ServiceWiring = wire.NewSet(
		llm.NewService,
	)
	UseCaseWiring = wire.NewSet(
		CommonWiring,
		RepositoryWiring,
		ServiceWiring,
		NotifierWiring,
		WSWiring,
		auth.NewUseCase,
		chat.NewUseCase,
		comment.NewUseCase,
		post.NewUseCase,
		like.NewUseCase,
		media.NewUseCase,
		notification.NewUseCase,
		profile.NewUseCase,
		llm.NewUseCase,
	)
	HandlerWiring = wire.NewSet(
		UseCaseWiring,
		auth.NewHandler,
		chat.NewHandler,
		comment.NewHandler,
		post.NewHandler,
		like.NewHandler,
		media.NewHandler,
		notification.NewHandler,
		profile.NewHandler,
		ws.NewBroadcaster,
		ws.NewHandler,
		llm.NewHandler,
	)
)

func InitServer(ctx context.Context) (server.Server, error) {
	panic(wire.Build(
		HandlerWiring,
		routes.NewRouter,
		server.NewServer,
	))
}
