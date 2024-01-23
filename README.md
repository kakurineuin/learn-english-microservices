# learn-english-microservices

## 介紹

這是使用 Golang 和 Go Kit、Echo、gRPC 開發的英文學習網站，使用 CircleCI 佈署到雲端主機。\
分成 3 個微服務：

- ExamService
  - 管理使用者建立的測驗、問題。
  - 提供測驗功能，記錄每次測驗結果。
  - 使用 gRPC 監聽 WebService 的請求。
  - 使用 Go Kit 架構。
- WordService
  - 查詢單字解釋
    - 使用 Golang 的 Colly 爬蟲程式去抓取線上英文字典網站的單字解釋。
  - 收藏最愛的單字解釋，產生單字卡以供複習。
  - 使用 gRPC 監聽 WebService 的請求。
  - 使用 Go Kit 架構。
- WebService
  - 運行網站
    - 前端使用 React、Redux、TypeScript
    - 後端使用 Golang、Echo(Golang 的 Web Framework)。
  - 處理使用者註冊、登入。
  - 單字卡功能使用 Redis 對 WordService 的回應結果作 cache
  - 使用 gRPC 呼叫 ExamService、WordService。

最初的想法是想以 K8S 的架構佈署在 Google Kubernetes Engine (GKE)，
但是感覺收費會很高，所以放棄這個選項。
改用 fly.io，在 fly.io 的 VM 裡面啟動這三個服務(不同的 port)，彼此使用 gRPC 通訊。

### K8S 的架構圖

![Alt text](image/%E7%B6%B2%E7%AB%99%E6%9C%8D%E5%8B%99%E6%9E%B6%E6%A7%8B_K8S.png?raw=true "網站服務架構_K8S.png")

### Fly.IO 的架構圖

![Alt text](image/%E7%B6%B2%E7%AB%99%E6%9C%8D%E5%8B%99%E6%9E%B6%E6%A7%8B_FlyIO.png?raw=true "網站服務架構_FlyIO.png")

## 使用的服務

資料庫：MongoDB, Redis\
CI：CircleCI\
雲端主機：fly.io\
網址：https://learn-english-microservices.fly.dev

## 網站使用說明

- 網站佈署在 fly.io，
  為了節省資源，設定了閒置時自動關閉的功能，
  所以第一次瀏覽或閒置一會後再瀏覽時，
  需要等個 4~5 秒。
- 登入才能使用每個功能。
- 註冊只需要帳號和密碼，不需要提供 Email。
- 若不想特地註冊一個帳號，可以使用以下帳號登入。
  - Username: guest01
  - Password: 12345678
- 每個功能的標題旁都會有小圓形驚嘆號，點擊它會出現功能說明。

![Alt text](image/home.png?raw=true "Home")

- 單字卡的操作方式
  - 滑鼠操作
    - 點擊卡片即可翻面
    - 點擊 \[下一張\] 移到下一張卡片
    - 點擊 \[上一張\] 移到上一張卡片
  - 鍵盤操作
    - 按 S 鍵即可翻面
    - 按 D 鍵移到下一張卡片
    - 按 A 鍵移到上一張卡片

![Alt text](image/word-card.png?raw=true "Word Card")

## 目錄與檔案說明

- .circleci
  - CircleCI config 的目錄
- ExamService
  - 提供測驗功能的微服務的目錄
- ProtoBuf
  - gRPC .proto 的目錄
- WebService
  - 運行網站的微服務的目錄
- WordService
  - 提供單字查詢與單字卡複習的微服務的目錄
- k8s
  - minikube 運行 k8s 的配置檔的目錄
- script
  - 佈署用的 .sh 的目錄
- skaffold.yaml
  - 本機使用 minikube 和 skaffold 測試 k8s 用的配置檔。
