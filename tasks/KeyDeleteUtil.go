package tasks

import (
	"fmt"
	"marlin/datamodels"
	"marlin/services"
	"time"
)

//DeleteKeys ,delete keys
func DeleteKeys(batchInfoService services.UserBatchInfoService, batchDetailService services.UserBatchDetailService,
	keyService services.KeyService) {
	go func(batchInfoService services.UserBatchInfoService, batchDetailService services.UserBatchDetailService,
		keyService services.KeyService) {
		for {

			result := batchInfoService.GetUnfinishBatch()
			batchs := result.Data.([]datamodels.UserBatchInfo)
			for _, batch := range batchs {
				fmt.Println(batch)
				detail := batchDetailService.GetDetailstByBatchID(int64(batch.ID))
				keys := detail.Data.([]datamodels.UserBatchDetail)
				for _, key := range keys {
					fmt.Printf("One key[%v]\n", key)
					rtn, err := keyService.DeleteByID(key.KeyName, "codis-demo")
					fmt.Printf("PostDel controller ,rtn:[%v],err:[%v]\n", rtn, err)
					var (
						msg  string
						code int = 0
					)
					if err == nil {
						if rtn == 1 {
							msg = "delete key ok."
						}
						if rtn == 0 {
							code = -2
							msg = "key not exists."
						}
					} else {
						code = -1
						msg = "delete key error." + err.Error()
					}
					fmt.Printf("Del key:[%s],code[%d],msg:[%s]\n", key.KeyName, code, msg)
					batchDetailService.UpdateByID(key)
				}
				batchInfoService.UpdateByID(batch)
			}
			fmt.Printf("del go routine running...")
			time.Sleep(time.Minute * 5)
		}
	}(batchInfoService, batchDetailService, keyService)

}
