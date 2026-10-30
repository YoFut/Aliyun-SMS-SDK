//go:build wireinject

package app

import (
	"log/slog"
	"net/http"

	"github.com/YoFut/Aliyun-SMS/internal/conf"
	"github.com/YoFut/Aliyun-SMS/internal/data"
	"github.com/YoFut/Aliyun-SMS/internal/web/api"
	"github.com/google/wire"
)

func WireApp(bc *conf.Bootstrap, log *slog.Logger) (http.Handler, func(), error) {
	panic(wire.Build(data.ProviderSet, api.ProviderVersionSet, api.ProviderSet))
}
