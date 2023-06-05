package tools

import "encoding/json"

type Metric struct {
	Keys      map[string]string      `json:"keys"`
	Vals      map[string]interface{} `json:"vals"`
	Timestamp string                 `json:"timestamp"`
}

type Metrics []Metric

func ProcessBytesMessage(msg []byte) (Metrics, error) {
	var m []Metric
	// 在这里处理字节类型的消息
	err := json.Unmarshal(msg, &m)
	if err != nil {
		return nil, err
	}
	// 可以根据需要进行解码、转换或其他处理操作
	return m, nil
}
