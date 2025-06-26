package global

import (
	"boost/data/server/config"
	"boost/data/server/utils/timer"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	FULL_NODE_CATFEE_ENABLE bool
	CATFEE_FULL_NODE        string
	X_CATFEE_TOKEN          string

	TRON_GRPC_NODE                  string
	MAX_WORKER_COUNT                int
	MAX_RECURSION_DEPTH             int
	MAX_TRX_COUNT                   int64
	DAY_LIMIT                       int
	IGNORE_TRANSFER_USDT_AMOUNT_MIN int
	TRANSFER_USDT_AMOUNT_MIN        int
	TRANSFER_TX_COUNT_MIN           int

	FROM_ADDRESS_PART_PREFIX int
	FROM_ADDRESS_PART_SUFFIX int
	TRANSFER_AMOUNT_MIN      int
	TRANSFER_AMOUNT_MAX      int
	BALANCE_RANGE_MIN        int
	BALANCE_RANGE_MAX        int

	TRONSCAN_KEYS            []string
	TRONGRID_KEYS            []string
	TRANSFER_TRX_AMOUNT_LIST []string
	TX_COUNT_RANGE_MIN       int
	TX_COUNT_RANGE_MAX       int
	AVG_AMOUNT               int
	LOCAL_DB_ENABLE          bool
	SUOJINSUANFA             int64
	BOT_ENABLE               bool
	TRXFAST_USERNAME         string
	TRON_MONITOR_ADDRESSES   string
	TRXFAST_PASSWORD         string
	MASTER_ADDRESS_PK        string
	MASTER_ADDRESS           string
	CONTRACT_ADDRESS         string
	TRON_FULL_NODE           string
	GVA_BOT                  string
	GVA_BOT_CHAT_ID          string
	GVA_DB                   *gorm.DB
	GVA_DB_B                 *gorm.DB
	GVA_DB_C                 *gorm.DB
	GVA_DB_D                 *gorm.DB
	GVA_DB_E                 *gorm.DB
	GVA_DB_F                 *gorm.DB
	GVA_DB_Local             *gorm.DB
	GVA_CONFIG               config.Server
	GVA_VP                   *viper.Viper
	// GVA_LOG    *oplogging.Logger
	GVA_LOG                 *zap.Logger
	GVA_Timer               timer.Timer = timer.NewTimerTask()
	GVA_Concurrency_Control             = &singleflight.Group{}
	GVA_ROUTERS             gin.RoutesInfo

	TRANSFER_TRX_AMOUNTS []string
)
