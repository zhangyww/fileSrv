/**
 * @Author: zhangyw
 * @Description:
 * @File:  Collector
 * @Date: 2021/5/20 16:19
 */

package collector

import "fileSrv/model"

type Collector interface {
	Collect(fileInfo *model.FileInfo)
}
