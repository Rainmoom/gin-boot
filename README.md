github.com/Rainmoom/gin-boot
# gin开发脚手架

### 使用示例
``` Go
import (
	"github.com/Rainmoom/gin-boot/pkg/conf"
	"github.com/Rainmoom/gin-boot/pkg/server"
	"github.com/Rainmoom/gin-boot/pkg/server/router"
	"github.com/Rainmoom/gin-boot/pkg/storage/cache"
	"github.com/Rainmoom/gin-boot/pkg/storage/db/postgres"
)

func main() {
    svr := &server.Server{
        Init: func() {
            //初始化数据库
            postgres.Init(conf.Cfg.Postgres)
            //初始化redis
            cache.InitRC(conf.Cfg.Redis)
        },
        Routers: []router.Router{
            //配置路由
        },
    }
    //设置配置文件地址并启动项目
    svr.Run("./config.yaml")
}

```