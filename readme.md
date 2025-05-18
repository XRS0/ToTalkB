ToTalkB
├─ auth
│  ├─ auth.go
│  ├─ cmd
│  │  └─ main.go
│  ├─ db
│  │  └─ postgres.go
│  ├─ go.mod
│  ├─ go.sum
│  ├─ middleware
│  │  └─ middleware.go
│  └─ pkg
│     └─ user.go
├─ chat
│  ├─ client.go
│  ├─ cmd
│  │  └─ main.go
│  ├─ go.mod
│  ├─ go.sum
│  ├─ home.html
│  ├─ hub.go
│  └─ pkg
│     ├─ chat.go
│     └─ message.go
├─ codes
│  ├─ cmd
│  │  └─ main.go
│  ├─ generate.go
│  ├─ go.mod
│  ├─ go.sum
│  ├─ hello.jpeg
│  └─ scan.go
├─ config
│  └─ config.yaml
├─ docker-compose.yml
├─ event_manager
│  ├─ config
│  │  └─ config.yaml
│  ├─ event_manager
│  ├─ go.mod
│  ├─ go.sum
│  ├─ internal
│  │  ├─ application
│  │  │  ├─ event_queue_service.go
│  │  │  └─ event_service.go
│  │  ├─ config
│  │  │  └─ config.go
│  │  ├─ domain
│  │  │  ├─ event.go
│  │  │  ├─ event_queue.go
│  │  │  └─ event_repository.go
│  │  └─ infrastructure
│  │     ├─ grpc
│  │     │  ├─ event_service.go
│  │     │  ├─ notification.go
│  │     │  └─ server.go
│  │     └─ persistence
│  │        ├─ memory
│  │        │  └─ event_repository.go
│  │        └─ postgres
│  │           ├─ event_queue_repository.go
│  │           ├─ migrations
│  │           │  ├─ 000001_init.up.sql
│  │           │  ├─ 000002_create_event_queues.down.sql
│  │           │  └─ 000002_create_event_queues.up.sql
│  │           └─ repository.go
│  └─ main.go
├─ gen
├─ LICENSE
├─ notify
│  ├─ config
│  │  └─ config.yaml
│  ├─ go.mod
│  ├─ go.sum
│  ├─ internal
│  │  ├─ application
│  │  │  └─ notification_service.go
│  │  ├─ config
│  │  │  └─ config.go
│  │  ├─ domain
│  │  │  ├─ gen
│  │  │  │  ├─ notification.pb.go
│  │  │  │  └─ notification_grpc.pb.go
│  │  │  └─ notification.go
│  │  ├─ handlers
│  │  │  └─ handler.go
│  │  ├─ infrastructure
│  │  │  ├─ grpc
│  │  │  │  ├─ client.go
│  │  │  │  └─ server.go
│  │  │  ├─ persistence
│  │  │  │  └─ postgres
│  │  │  │     ├─ migrations
│  │  │  │     │  ├─ 000001_create_notifications_table.down.sql
│  │  │  │     │  └─ 000001_create_notifications_table.up.sql
│  │  │  │     └─ repository.go
│  │  │  └─ repository
│  │  │     └─ postgres
│  │  │        └─ notification.go
│  │  ├─ server
│  │  │  ├─ http.go
│  │  │  └─ server.go
│  │  ├─ service
│  │  │  └─ notification.go
│  │  └─ websocket
│  │     ├─ handler.go
│  │     ├─ manager.go
│  │     └─ notification_handler.go
│  ├─ main.go
│  └─ notify
├─ proto
│  ├─ event.proto
│  ├─ gen_event
│  │  ├─ event.pb.go
│  │  ├─ event_grpc.pb.go
│  │  ├─ go.mod
│  │  └─ go.sum
│  ├─ gen_notify
│  │  ├─ go.mod
│  │  ├─ go.sum
│  │  ├─ notification.pb.go
│  │  └─ notification_grpc.pb.go
│  └─ notification.proto
├─ readme.md
├─ schema
│  ├─ 000001_init.down.sql
│  └─ 000001_init.up.sql
└─ test_client
   ├─ client
   │  └─ test_client.go
   ├─ cmd
   │  └─ main.go
   ├─ gen
   │  ├─ event.pb.go
   │  └─ event_grpc.pb.go
   ├─ go.mod
   ├─ go.sum
   └─ main.go
