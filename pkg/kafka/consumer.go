package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"

	"little-vote/pkg/dao"
	"little-vote/pkg/viper"
)

func StartConsumer(ctx context.Context) {
	consumer, err := sarama.NewConsumer([]string{viper.Config.GetString("kafka.address")}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v\n", err)
		return
	}
	partitionList, err := consumer.Partitions("vote")
	if err != nil {
		fmt.Printf("fail to get list of partition:err%v\n", err)
		return
	}
	fmt.Printf("partition list:%v\n", partitionList)
	fmt.Println("start consumer success")
	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("vote", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		defer pc.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				err := dao.SyncUser(string(msg.Value))
				if err != nil {
					fmt.Println(err)
				}
			}
		}(pc)
	}
	<-ctx.Done()
}
