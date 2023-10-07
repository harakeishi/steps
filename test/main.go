package main

import (
	steps "github.com/harakeishi/steps"
)

type TargetData struct {
	TargetName string
	TargetAge  int
}

func main() {
	// 新規フローを作成
	flow := steps.NewFlow()
	// ステップを追加
	flow.AddStep(
		steps.NewStep(
			"Get Target Data",
			"処理の対象者を取得する",
			func(inputs interface{}) (interface{}, error) {
				// 処理の対象者を取得する処理
				targetData := []TargetData{
					{
						TargetName: "Taro",
						TargetAge:  20,
					},
					{
						TargetName: "Jiro",
						TargetAge:  30,
					},
				}
				return targetData, nil
			},
			0,
			&steps.Step{},
		))
	flow.AddStep(
		steps.NewStep(
			"Print Target Data",
			"処理の対象者を表示する",
			func(inputs interface{}) (interface{}, error) {
				// 処理の対象者を表示する処理
				targetData := inputs.([]TargetData)
				for _, data := range targetData {
					println(data.TargetName)
					println(data.TargetAge)
				}
				return nil, nil
			},
			0,
			flow.Steps[0],
		))
	// フローを実行
	flow.Run()
}
