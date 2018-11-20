# 前置き

* [Kubernetes完全ガイド](http://amzn.asia/d/9QFJw5X)を読むのが一番早い。分かりやすく網羅的。

# Kubernetesについて

## Kubernetesとは何か

* 読み方
    * 「クーバネィテス」と呼ぶ人が多いようです。  https://twitter.com/IanMLewis/status/978831106786013185
        * 日本人は「クバネテス」「クバネティス」って呼んでいる人多い気がする。
    * 好きに呼べばいいと思います。
*  Kubernetesは、**コンテナオーケストレーションツールと呼ばれるアプリケーションです**。コンテナ化されたアプリケーションのデプロイ、スケーリングなどの管理を自動化することができます。
* コンテナオーケストレーションツールにはDocker Swarm、 Amazon ECS、 Rancher、 Apache Mesos、 Nomadなどがありますが、Kubernetesが今ではコンテナオーケストレーションツールのデファクトスタンダードとなりました。
* [DockerがKubernetesとの統合およびサポートを発表。DockerCon EU 2017](http://www.publickey1.jp/blog/17/dockerkubernetesdockercon_eu_2017.html) 「Kubernetesがオーケストレータにおける事実上の標準の地位を固めたと言えそうです。」
* [「Kubernetes is becoming the Linux of the cloud」 - Jim Zemlin, Linux Foundation](https://twitter.com/kubernetesio/status/840257886202683392)


## 歴史

* KubernetesはGoogleが社内で利用されていたコンテナクラスタマネージャの「Borg」と呼ばれるシステムが元になっている。
* 開発、管理は現在、[Cloud Native Computing Foundation(CNCF)]( https://www.cncf.io/)に移管されています。CNCFには著名な開発者、エンドユーザ、大手クラウドプロバイダなどのベンダーが参加しており、現在はCNCFが主体となり中立的な立場で開発が進められています。
* CNCFではKubernetesの仕様を標準化していっており、クラウドプロバイダ間で同じように使うことができる。つまりGCPのKubernetes上で動かしているアプリケーションをAWSのKubernetesで同じように動かすことができる。これによりクラウドベンダーロックインを避けることが可能となっている。
* 2014/11にGCPがマネージドKubernetesサービスとしてGKEをリリース。2017/2にAzureがAKS(Azure Kubernetes Service)をリリース。2017/11にAWSがAmazon EKS(Amazon Elastic Container Service for Kubernetes)をリリース。
* 2014年に発表、2015にv1.0、2016/12にv1.5、2017/12にv1.9。非常に進化が速い。が、今年から落ち着いてきた。


## Kubernetesを使うと何が出来るのか

* コンテナのデプロイ
    * 設定ファイルをアップすればそれに基づいてデプロイしてくれる
    * 特別な指定がない場合には、CPUやメモリの空きリソースの状況に従ってデプロイが行われるため、ユーザはどのサーバーに配置するかなどを管理する必要がありません
    * もちろん特定のサーバーに特定のアプリケーションをデプロイする機能もある。このアプリケーションはSSDのサーバーにデプロイするとか
* サービスを停止することなくデプロイが可能(ローリングアップデート)
    * またデプロイのバージョン管理ができる。つまり1つ前にロールバックしたり、指定したバージョンに戻すことが可能
* オートスケーリング
    * 負荷に応じて自動的にスケールさせることができる
    * クラウドプロバイダによっては負荷に応じてサーバーを追加してくれる
    * AWSのオートスケーリングをイメージしてもらえれば
* オートヒーリング(ヘルスチェック、自動復旧)
    * ヘルスチェック機能があり、ヘルスチェックに失敗したら自動で新しくデプロイしてくれる。ヘルスチェックの内容はHTTP/TCP・シェルスクリプトで設定できる
    * クラスタのサーバー(ノード)に障害が起きた場合でも、自動で他のサーバー(Node)に再デプロイしてくれる
    * コンテナだけではなくサーバー(Node)自体も自動復旧対象
* ロードバランシング
* サービスディスカバリ
* データの管理
    * Kubernetesクラスタの情報の保存にはetcdが使われ、冗長化され保存される。コンテナが使用する認証情報や設定情報もetcdに保存される。
* 外部ストレージの使用
* ログの管理
* 権限管理
* Infrastructure as Code
    * KubernetesはYAML(or json)で書かれたファイルで、デプロイするコンテナや周辺リソースを管理します
* その他エコシステムとの連携や拡張



# Kubernetesクラスタの構成

`Master`と`Node`に分かれています。
`Master`と`Node`を合わせて`Cluster`と呼ぶ。

* Master
    * **Kubernetesシステム全体を管理するサーバー**
        * 機能としては以下
            * Kubernetes API Server
                * APIを提供する
            * Scheduler
                * Nodeにアプリケーションをデプロイする機能
                * ちなみに、Kubernetesではデプロイすることをスケジュールするといいます
            * Controller Manager
                * アプリケーションのヘルスチェックや自動復旧を行う機能
            * etcd
                * クラスタ設定を保存する高可用分散KVS
        * Masterの冗長構成も可能
* Node
    * **デプロイされたアプリケーションが実行されるサーバー**
        * 機能としては以下
            * Docker等のcontainer runtime
            * Kubelet
                * MasterのAPI Serverとの通信
                * containerの管理
            * Kubernetes Service Proxy (kube-proxy)
                * アプリケーション間のネットワークトラフィックのロードバランシング

# Kubernetesの各リソース

色んなリソースがありますが、まずは、Pod, ReplicaSet, Deployment, Service, IngressをおさえておけばOKです。

## Pod

* Kubernetes上でのデプロイの最小の単位。Pod単位でデプロイされる。
* Podには複数のコンテナを入れることができる。
* Podがデプロイ単位なので、pod内のコンテナたちはまとめてデプロイされ、同じnodeにデプロイされる。
* なのでWebアプリケーションが入ったコンテナとDBコンテナは同じPodに入れないほうがいい。一緒のサーバーにデプロイされてしまうため。通常WebサーバーとDBサーバーは分けますよね？
* 1podにつき1コンテナになることが多いが、関連性が強いものは1podに複数コンテナ入れる。
    * たとえばWebアプリケーションコンテナとログ収集コンテナを1つのPodに入れるとか
* pod内は論理的な1つのホストとして扱われ、pod内の各コンテナは以下を共有する
    * PIDネームスペース （pod内のアプリケーション群からはお互いのプロセスが見える）
    * ネットワーク・ネームスペース （pod内のアプリケーション群は同一のIPとポートの空間にアクセスする）
    * IPCネームスペース （pod内のアプリケーション群は、SystemV IPCまたはPOSIXメッセージキューを使って相互通信できる）
    * UTSネームスペース （pod内のアプリケーション群はホスト名を共有する）
* Podにはクラスタ内で有効な内部IPアドレス(ClusterIP)が付与される。
* Pod単体では外部のネットワークからアクセスができない。外部のネットワークからPodへアクセスするためには後述するServiceを使用する必要がある。

![スクリーンショット 2018-11-18 14.35.10.png](https://qiita-image-store.s3.amazonaws.com/0/14124/b25ca724-f1bf-670a-1a50-19a0c9e59f1d.png)
NodeをまたがってPodはデプロイされない。

## ReplicaSet

* Podのスケール制御をするリソース
* 設定で定義された数のpodを作成、維持します。
    * 例えばpodが削除された等で指定されたpod数に満たなくなった場合、ReplicaSetは数を満たすようにpodを作成します。
* 後述するDeploymentを使うことを推奨されているため、ユーザがこのリソースを単体で使うことはほぼない。

## Deployment

* ReplicaSetと同じくpodのスケーリング制御ができるが、さらにデプロイの世代管理やローリングアップデートもできる。
* **Deploymentを作成すると、ReplicaSetとPodが自動で生成される。**
* 世代管理というのは、新しいPodをデプロイしたときに古いPodは削除せずに取っておいてバグ等の何か問題があったときに古いバージョンのPodに戻せる機能。
* ローリングアップデートというのは、1つPodを停止状態にし、1つ新しくPodをデプロイすることを繰り返すことでシステムを無停止でデプロイできるというもの。


## Service

* Serviceは複数のPodを1つにまとめたエンドポイントを提供するためのリソースです
* まずPodは増えたり減ったり再作成されたりするので内部IPが一定ではありません。またそもそも外部からアクセスができません。これらを解決するためにServiceがあります。たとえばWebサーバー用Podが複数あったとして、それらのPodに1つのIPアドレスや(内部的な)ドメインを通してアクセスしたい場合はServiceを使用すると可能になります。
    * VIPを想像すると理解しやすいかと思います。
* ServiceはL4ロードバランサの役割を担当する
    * L7ロードバランサは後述するIngressが担当する
* Serviceにはエンドポイントを提供する複数の種類がありますので、要件に合わせて選択しましよう。



### Serviceの種類

|名称|説明|
|:--|:--|
|ClusterIP|基本となるタイプ。クラスタ内部で有効なDNS名/仮想IPで複数のPodを紐付ける。仮想IPはClusterIPと呼ばれる。|
|NodePort|ClusterIP Serviceを作成し、このServiceをNode上の静的なポート番号で公開する|
|LoadBalancer|クラスタ外からトラフィックを受ける場合に使い勝手がいいService。NodePort Serviceを作成し、プロバイダ特有の実装でロードバランサを作成して同Serviceと紐付ける。プロバイダによってはSSL処理を行うこともできるが、より汎用的に外部からのアクセスを処理するには後述するIngressを使う|
|ExternalName|外部のサービスを利用したいときに、DNSのCNAMEレコードを利用してクラスタ内部向けのエイリアスを作成する|
|Headless|Service名単体だけではDNS解決ができないが、紐付けられたPod名をDNSのSRVレコードに持つ。StatefulSetなどを利用するときに使う|

## Ingress

* IngressはServiceに紐付き、リクエストをServiceにプロキシします。簡単にいうとIngressはL7ロードバランサです。
* アプリケーション層(L7)でのURLパターンによる転送ルールの適用や、SSL終端処理などもIngressが行います。
* IngressはKubernetesが動作している環境によって機能とその実装が大きく変わり、GCPではGCPの、AWSではAWSの環境に合った形でロードバランサーが作成されます。
    * たとえばGCP上ではGCPのHTTPロードバランサが作成されます。
* GCP上だとIngressリソースを定義すると自動的にHTTPロードバランサが作成されるが、オンプレだと作成されないため自前でロードバランサを用意する必要がある。。


## DaemonSet

* 各Nodeに1つだけpodを配置するためのリソース
* たとえばNodeのメトリックス(CPU使用率とか)を取得するためのPodはNodeに1つだけあれば十分なのでそういった用途にはDaemonSetが使われる。
* 特定のNodeだけにPodを配置することもできる

## Job/CronJob

* Jobは一回きりのタスクを実行するためのリソース
* CronJobはバッチのように定期的に実行するためのリソース


## Volume

* ストレージを表すリソース
* KubernetesにおけるVolumeはDockerにおけるそれとほぼ同じで、何らかの外部ストレージ領域をVolumeとしてコンテナにマウントできます。
* Volumeにはさまざまな種類があり、クラウドプロバイダ依存のgcePersistentDiskやAzureDiskVolumeやawsElasticBlockStore、自前でミドルウェアやハードウェアを用意するならばNFSやiSCSI等が使えます。

## ConfigMap/Secret

* 設定値を保持するためのリソース
* キーと値をセットで保持できる。
* コンテナイメージを作成する際、イメージそのものにアプリケーションの設定を書き込むと、設定を1つ変えるたびに新しいイメージを作成しなければなりません。そうはならないように、各種設定とコンテナに展開されるプログラムは切り離してデプロイできるようにする必要があります。コンテナの設定をPodとコンテナから切り離すために、KubernetesにはConfigMap、Secretというリソースがあります。ConfigMapもSecretもそれぞれ、キーと任意の値を登録するものです。
* ConfigMapでは、プレーンテキストとして格納されます。
* Secretでは一応エンコードされますが、値はbase64エンコードされているだけです。Kubernetesにアクセス可能なユーザーなら誰でも中身を見ることができますので、Secretという名前から安全と考えてはいけません。Secretが定義されたマニフェストファイルを暗号化する[kubesec](https://github.com/shyiko/kubesec)というものを使う選択肢がある。

## その他

### Label

* Kubernetesではいろんなリソース(podやserviceやdeployment等々)にLabelをつけることができる。つまりリソースのグルーピングができる。
* Labelは名前と値で構成されている。
* 何か操作するときに、Label情報をもとに操作することができる。
* たとえば複数podをデプロイするときに、そのpodに同じLabelをつけておけば、複数のpodをLabel情報をもとに操作が可能になる。
* AWSのタグと似ている。

# Kubernetesをサポートしている外部システム

- Ansible
    - Kubernetesへのコンテナデプロイ
- Fluentd
    - Kubernetes上のコンテナのログを転送
- Jenkins
    - Kubernetesへのコンテナのデプロイ
- [Jenkins X](https://www.publickey1.jp/blog/18/jenkins_xgitdockerkubernetescicd.html)
    - CIをKubernetesクラスタ上で動かせるJenkins
- [Istio](https://istio.io/)
    - Kubernetes上にサービスメッシュを構築 
- [Prometheus](https://prometheus.io/)
    - Kubernetesのモニタリング
- [Spinnaker](https://www.spinnaker.io/)
    - Kubernetesへのコンテナデプロイ
- [Kubeflow](https://www.kubeflow.org/)
    - Kubernetes上にML Platformをデプロイ
- [Argo](https://argoproj.github.io/)
    - Kubernetes上で動作するワークフローエンジン
- [Rook](https://rook.io/)
    - Kubernetes上に分散ファイルシステムを構築
- [Vitess](https://vitess.io/)
    - Kubernetes上にMySQLクラスタを構築
- Spark
    - jobをKubernetes上でネイティブ実行
