/*
 * @Author: David
 * @Date: 2020-03-08 21:01:58
 * @LastEditTime: 2020-03-10 18:13:15
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/datamodels/DbResult.go
 */

package datamodels

// DbResult db result
type DbResult struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewDbResult new db result
func NewDbResult(data interface{}, c int, m ...string) *DbResult {
	r := &DbResult{Data: data, Code: c}

	if e, ok := data.(error); ok {
		if m == nil {
			r.Msg = e.Error()
		}
	} else {
		r.Msg = "SUCCESS"
	}
	if len(m) > 0 {
		r.Msg = m[0]
	}

	return r
}
