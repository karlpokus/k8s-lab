# app
- gw: the gateway. Requires auth via the user service. Calls blog service if successful.
- user: r/w the users db
- blog: r/w the blog db

# todos
- [x] GET /posts
- [x] GET /post/:title
- [x] POST /post
- [ ] api tests
- [ ] run test on merge
- [ ] readiness and liveness probes
- [ ] add binary version
- [x] add graceful exits to gw
- [x] add request logs toggle from env
- [ ] create configmap from file
- [x] elastic apm tracing
- [ ] logging
- [ ] 4 golden sigs
- [ ] multi-stage docker builds
- [ ] ingress

# apm
Can elastic apm provide 4 golden signals?

- traffic: if ELASTIC_APM_TRANSACTION_SAMPLE_RATE <1 then we don't have all data
- latency: for gw at least. service.name: gw, transaction.duration.us, transaction.name, processor.event: "transaction", http.response.status_code OR transaction.result. Mongo latency span.subtype:mongodb and span.duration.us
- errors: 5xx/s
- saturation (per service); cpu: system.process.cpu.total.norm.pct, mem: system.process.memory.rss.bytes

Availability (not part of 4 sigs) could be calculated by (reqs - errors) / reqs

# usage
We'll use docker-compose to orchestrate containers for local development and we'll keep images simple i.e huge - not optimized, to cut down on rebuild time.

```bash
# run all
$ docker-compose -f mongo.yml -f apps.yml up &
# rebuild image and restart container on src file change
$ watchexec -w src/gw "docker-compose -f apps.yml up -d --build gw" & \
watchexec -w src/user "docker-compose -f apps.yml up -d --build user" & \
watchexec -w src/blog "docker-compose -f apps.yml up -d --build blog" &
# stop
$ test `jobs -p | wc -l` -gt 0 && kill `jobs -p`
```

known bugs
- docker-compose service names change colour between restarts
- go build errors end up in docker process list

# build
```bash
# tag
$ docker tag app_gw:latest pokus2000/gw:0.3.0
$ docker tag app_user:latest pokus2000/user:0.3.0
$ docker tag app_blog:latest pokus2000/blog:0.4.0
# push
$ docker push <imgs>
# note: update image tag in deployment
```

# deploy
```bash
# update cluster
$ kubectl apply -f app/src -R && kubectl get pods -w | grep -i running
# note: if only configmap change: kubectl rollout restart deploy/<thing>
```
