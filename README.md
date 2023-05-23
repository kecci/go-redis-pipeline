# go-redis-pipeline

## Usecase
- Set & Del 1000 Keys.
- Set & Del 1 Keys.

## Functions 
1. setCommand - Set keys with single command.
2. setCommandPipeline - Set keys with pipeline.
3. delCommand - Delete keys with single command.
4. delCommandPipeline - Delete keys with pipeline.

## Output
```
# 1000 keys
2023/05/23 10:42:03 processing 1000 keys time delCommandPipeline: 1.132291ms or 1.132µs/key
2023/05/23 10:42:03 processing 1000 keys time delCommand: 887.208µs or 887ns/key
2023/05/23 10:42:03 processing 1000 keys time setCommandPipeline: 4.451167ms or 4.451µs/key
2023/05/23 10:42:03 processing 1000 keys time setCommand: 606.724667ms or 606.724µs/key

# 1 key
2023/05/23 10:45:31 processing 1 keys time delCommandPipeline: 569.542µs or 569.542µs/key
2023/05/23 10:45:31 processing 1 keys time delCommand: 532µs or 532µs/key
2023/05/23 10:45:31 processing 1 keys time setCommandPipeline: 574.833µs or 574.833µs/key
2023/05/23 10:45:31 processing 1 keys time setCommand: 12.914542ms or 12.914542ms/key
```

## Conclusion
We recommend using a pipeline to execute multiple commands, but we do not recommend using a pipeline if there is only one command because the result may be longer than a normal single command.

## Reference
- https://medium.com/@babagusandrian/redis-pipeline-vs-single-proccess-di-golang-6e9a0fb1fcb3