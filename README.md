# go-redis-pipeline

## Usecase
- Set 10.000 Keys.
- Set 1 Keys.

## Functions 
1. singleCommandRedis - Set 10.000 keys using a single command per loop.
2. pipelineRedis - Set 10.000 keys using a pipeline.
3. singlePipelineRedis - Set 1 key using a pipeline.

## Output
```
2023/05/19 23:53:30 processing 10000 keys time singleCommandRedis: 4.64399775s or 464.399µs/key
2023/05/19 23:53:30 processing 10000 keys time pipelineRedis: 26.64325ms or 2.664µs/key
2023/05/19 23:53:30 processing 1 keys time singlePipelineRedis: 472.917µs or 472.917µs/key
```

## Summary
- Single 10.000 Keys -> 464.399µs/key
- Pipeline 10.000 Keys -> 2.664µs/key
- Pipeline 1 Key -> 472.917µs/key


## Conclusion
We recommend using a pipeline to execute multiple commands, but we do not recommend using a pipeline if there is only one command because the result may be longer than a normal single command.

## Reference
- https://medium.com/@babagusandrian/redis-pipeline-vs-single-proccess-di-golang-6e9a0fb1fcb3