
# 0, 用語の理解

Kubernetesは概念が命。
用語を理解しないと始まらない。 (そして用語が沢山あるのだ…)

[Kubernetesの用語の理解](https://qiita.com/sheepland/private/cdff472ba2d37784a125)


# 1, 準備

### 事前準備

- ローカルにpython 2.7がインストールされていことを確認
- `brew install wget`

### GCPのアカウントを作成

会社のメアドではなく、個人のメアドでアカウントを作成します。
会社のメアドで作るとstudyplus会社のGCPアカウントの中でユーザ作成がされて、同名のプロジェクトが作れなくなって面倒なことになります。

### GCPのコンソールにアクセスできることを確認。

https://console.cloud.google.com/home/dashboard

### プロジェクトの作成

この[リソースの管理](https://console.cloud.google.com/cloud-resource-manager)画面からプロジェクトを作成してください。
プロジェクト名は`k8s-test`にします。

### よく使う画面をピン留めする

以下の4つの画面はこのハンズオンでよく見るのでピン留めしておきます。(左上のメニューからサービス名にマウスオーバーするとピンのアイコンがでるのでクリックするとピン留めできます。ピン留めするとメニューの上のほうに固定されます)

![スクリーンショット 2018-10-15 22.33.23.png](https://qiita-image-store.s3.amazonaws.com/0/14124/bfb473eb-1ee7-d6eb-475b-3b0315f13301.png)

### gcloudツールのインストール

https://cloud.google.com/sdk/docs/quickstart-mac-os-x
の"始める前に"の1〜4を行います。
"始める前に"のセクションだけやればいいです。"SDK の初期化"のセクションはやらなくて大丈夫です。


### パスを通す
以下の設定で勝手にパスが通ります。

bashの人

.bashrc

```sh
# bashの人は以下。ただし"$HOME/path/to"の部分は適切に修正すること。
source $HOME/path/to/google-cloud-sdk/completion.bash.inc
source $HOME/path/to/google-cloud-sdk/path.bash.inc

# エイリアスはお好みで(kubectlコマンドは後でインストールするが先にaliasだけ設定しておく)
alias kc=kubectl 
```

zshの人

.zshrc

```sh
# zshの人は以下。ただし"$HOME/path/to"の部分は適切に修正すること。
source $HOME/path/to/google-cloud-sdk/completion.zsh.inc
source $HOME/path/to/google-cloud-sdk/path.zsh.inc

# エイリアスはお好みで(kubectlコマンドは後でインストールするが先にaliasだけ設定しておく)
alias kc=kubectl
```

```
# shellを再起動する
$ exec $SHELL -l
```

```
# gcloudコマンドが通ることを確認
$ gcloud version

# 最新版にアップデート(ちょっと時間がかかるかも)
$ gcloud components update
```

### プロジェクトIDの設定をする

```bash
$ PROJECT_ID=hogehoge #ここにプロジェクトIDを設定する(プロジェクトIDの取得の仕方は下の画像を参照)
```

[リソースの管理](https://console.cloud.google.com/cloud-resource-manager)画面に行き、以下をコピーする。
<img width="736" alt="スクリーンショット 2019-07-13 16.52.46.png" src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/14124/6d2fa4cf-c422-b9b2-02ec-a5b2e3fab039.png">



```bash
$ echo $PROJECT_ID # k8s-test以外の文字が出力されることを確認
$ gcloud config set project $PROJECT_ID
$ gcloud config list | grep project # 設定されていることを確認
```

### リージョンの設定をする

```bash
$ gcloud config set compute/zone asia-northeast1-a
```

### 認証をする

```bash
$ gcloud auth login
```

上記コマンドを入力するとブラウザが立ち上がりOAuth認証が始まりますので、プロジェクト作成に利用したアカウントで認証を行ってください。認証が成功すると、gcloudコマンドからGCPを操作できるようになります。

### kubectlのインストール

kubectlはKubernetesを操作するためのCLIツールです

```bash
$ gcloud components install kubectl
```

###  Google Cloud Platform サービスの有効化  

GCPでは各サービスを利用する際にAPIを有効化する必要があります。
GCPの管理画面から各サービスの画面にアクセスすると、自動でAPIが有効になります。
しかしAPIの有効化に1分弱かかるので、今回使うサービスを事前に有効化しておきます(実行完了まで少々かかります)。

```bash
$ gcloud services enable container.googleapis.com compute.googleapis.com cloudbuild.googleapis.com containerregistry.googleapis.com
```

# 2, ハンズオンの準備

```bash
$ cd 適当な作業ディレクトリ
$ git clone https://github.com/akira-kuriyama/kubernetes-handson.git
$ cd kubernetes-handson/source
$ echo $PROJECT_ID # k8s-test以外の文字が出力されることを確認
$ sed -i '' "s/PROJECT_ID/$PROJECT_ID/g" *.yaml
```

以降は `kubernetes-handson/source`以下のyamlファイルを使ってハンズオンを進めます。

# 3, Kubernetes Clusterの作成

Kubernetes Clusterの作成します。(デフォルトで3台のNodeが作成されます)

```bash
$ gcloud container clusters create my-k8s
```

https://console.cloud.google.com/kubernetes で確認できます。
もしくは`$ gcloud container clusters list`で確認できる。

完了まで数分(5,6分)かかる。

Kubernetes Clusterが作成されると以下のようになる。
<img width="1084" alt="スクリーンショット 2019-07-13 16.57.00.png" src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/14124/298da12e-b5bd-e086-da17-f6a8ecac0c27.png">


また、 https://console.cloud.google.com/compute/instances をみるとGCEインスタンスが作成されているのが分かる。これがnodeとなる。

クラスタが作成できたら、kubectlが正しくクラスタに接続して操作を行えるように、kubernetesクラスタの認証情報をセットします。

```bash
$ gcloud container clusters get-credentials my-k8s
``` 

確認

```bash
$ kubectl cluster-info
```

以下のような表示がされたらOK

```bash
$ kubectl cluster-info
Kubernetes master is running at https://35.200.6.30
GLBCDefaultBackend is running at https://35.200.6.30/api/v1/namespaces/kube-system/services/default-http-backend/proxy
Heapster is running at https://35.200.6.30/api/v1/namespaces/kube-system/services/heapster/proxy
KubeDNS is running at https://35.200.6.30/api/v1/namespaces/kube-system/services/kube-dns/proxy
kubernetes-dashboard is running at https://35.200.6.30/api/v1/namespaces/kube-system/services/kubernetes-dashboard/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```



# 4, podを配置

これから"◯◯ファイルを作成"というのが出てきますが、すでにファイルはgit cloneして手に入れているのでスキップして下さい。

### 最初にimageを作成

main.goを作成。8080ポートでまちうけ、Hello, World!を返すプログラム。


hello.go

```go
package main

import (
        "io"
        "net/http"
        "os"
)

func main() {
        http.HandleFunc("/", hello)
        http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        // 出力メッセージは環境変数があればそれを使い、なければ"Hello, World!"を使う
        message := os.Getenv("CUSTOM_MESSAGE")
        if message == "" {
                message = "Hello, World!"
        }
        io.WriteString(w, message)
}
```

次にDockerfileを作成

Dockerfile

```txt
FROM alpine:3.6
EXPOSE 8080
ADD hello-world /hello-world 
CMD ["/hello-world"]
```

次にcloudbuild.yamlを作成

cloudbuild.yaml

```yaml
steps:
- name: 'gcr.io/cloud-builders/go:alpine'
  env: ['PROJECT_ROOT=my-project']
  args: ['build', '-o', 'hello-world', 'hello.go']
- name: 'gcr.io/cloud-builders/docker'
  env: ['PROJECT_ROOT=my-project']
  args: ['build', '--tag=asia.gcr.io/PROJECT_ID/my-project/hello-world:latest', '.']
images: ['asia.gcr.io/PROJECT_ID/my-project/hello-world:latest']
```

次にGoogle Container Register(以下、GCR) にイメージを登録します。

```
$ gcloud builds submit --config=cloudbuild.yaml .
```

https://console.cloud.google.com/gcr にアクセスして登録されたことを確認。

`$ gcloud builds list`でも確認できる。

### Podの作成

次にpodファイルを作成

pod.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello-world
spec:
  containers:
    - image: asia.gcr.io/PROJECT_ID/my-project/hello-world
      imagePullPolicy: Always
      name: hello-world
```

そして以下のコマンドを実行して、podをKubernetes上に作成

```
kubectl create -f pod.yaml
```

podが作成されたか確認する

```bash
$ kubectl get pods
NAME          READY     STATUS    RESTARTS   AGE
hello-world   1/1       Running   0          5s
```

※ ちなみに`kubectl get pod`とか `kubectl get po`でもいける。

またGCPのKubernetes画面のワークロードの画面でも確認できる。
https://console.cloud.google.com/kubernetes/workload

pod内のコンテナにコマンドを実行する。

```
kubectl exec -it pod名(kubectl get podの左側のNAME) date
```

pod内のコンテナにログインする

```
kubectl exec -it pod名(kubectl get podの左側のNAME) /bin/sh
```

pod内に複数のコンテナが動いている場合は、以下のようにコンテナ名を指定する

```
kubectl exec $pod_name --container $container_name $command
```

また、
`$ kubectl get pod -o wide`で追加情報(どのnodeで動いているか等)が表示される。

get pod以外にも

```
kubectl get deployments
kubectl get replicasets
kubectl get services
```

があり、それぞれ `-o wide`が使える。

あとpod(やserviceやreplicasetやdeployment)の定義を知りたい場合は、
`kubectl get pod -o yaml`のように`-o yaml`をつけるとみれる。便利。

podの詳細情報を見たい場合は、
`kubectl describe pod`や`kubectl describe pod pod名`でみれる

コンテナが起動しない場合は、`kubectl describe pod pod名`の`Events`欄をみると原因が分かることが多い。

# 5, podへのアクセス

pod内のコンテナへ簡単にアクセスするための方法として、port forwardingがあります。

```
$ kubectl port-forward pod名(kubectl get podの左側のNAME) 8080
```

で、別タブで、以下を実行すると`Hello, World!`が返ってきます。

```
$ curl http://localhost:8080/
```


確認できたら、いったん削除しましょう

```
$ kubectl delete pods pod名(kubectl get podの左側のNAME)
```

もしくは

```sh
# こっちのほうがよく使う
$ kubectl delete -f pod.yaml
``` 

# 6, 設定値を設定ファイルから与える

まず以下のような設定ファイルを作成

configmap.yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: hello-world-config
data:
  hello-world.message: "Hello World, ConfigMap!"
```

以下を実行する。

```sh
$ kubectl create -f configmap.yaml
# 以下で作成されていることを確認
$ kubectl get configmaps
$ kubectl describe configmaps
```

次にpodリソースファイルを作成。envセクションがポイント。configMapKeyRefの部分で設定ファイル名と使用する設定Keyを指定している。

pod-configmap.yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: hello-world-configmap
spec:
  containers:
    - image: asia.gcr.io/PROJECT_ID/my-project/hello-world
      imagePullPolicy: Always
      name: hello-world
      env:
        - name: CUSTOM_MESSAGE
          valueFrom:
            configMapKeyRef:
              name: hello-world-config
              key: hello-world.message
```

で、
```
$ kubectl create -f pod-configmap.yaml
```

でpodを作成し、

```
$ kubectl port-forward pod名(kubectl get podの左側のNAME) 8080
``` 

をして、
別タブで、

```
$ curl http://localhost:8080/
```

すると、"Hello World, ConfigMap!"とでます

確認できたら、いったん削除しましょう

```
kubectl delete pods pod名(kubectl get podの左側のNAME)
```

もしくは

```
kubectl delete -f pod-configmap.yaml
```

# 7, Deploymentを使ったデプロイ

deployment.yaml

```yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: hello-world-deployment
spec:
  replicas: 3
  template:
    metadata:
      labels:
        name: hello-world-pod
    spec:
      containers:
        - image: asia.gcr.io/PROJECT_ID/my-project/hello-world
          imagePullPolicy: Always
          name: hello-world-pod
```

次にデプロイをする

```
$ kubectl create -f deployment.yaml
```

**これでdeploymentとreplicasetとpodが作成される。**

podはコンテナを擁するリソース、
replicasetはpodの数を維持するためのリソース、
deploymentはreplicasetを擁し、世代管理(バージョニング)をするためのリソース になります。

なので、deploymentの中にreplicasetがあり、replicasetの中にpodがあり、podの中にコンテナがあるイメージです。


```
$ kubectl get pod,replicasets,deployments
```

で確認

もしくは 

```
$ kubectl get all
```

で確認

GCPのKubernetesのワークロードの画面で、deploymentを確認することができる。
https://console.cloud.google.com/kubernetes/workload


```
$ kubectl get pods
```

を実行するとpodが3つ作成されているのが分かる。
試しにpodの一つを手動で消してみる。

```
$ kubectl delete pods pod名(kubectl get podの左側のNAME)
```

すると、podが3つに戻ることが確認できます。

```
$ kubectl get pods
```

つぎに、deployment.yaml内の`replicas: 3`を`replicas: 1`にして、
`$ kubectl replace -f replicaset.yaml` とやって設定を変更し、
`$ kubectl get pods`とすると、pod数が1になっていることが確認できます(すぐにpodが復活するで分かりづらいかもしれませんが、AGEの欄をみると若いpodが1ついるのが分かると思います)。

※ このdeploymentは消さなくてよいです。

# 8, Serviceの作成
 

LoadBalancerも一緒に作成してくれるServiceを作成します(`type: LoadBalancer`)。
このtypeはIngressがなくても静的IPが割り当てられます。

yaml内の`selector`に`name: hello-world-pod`と指定されているが、これは`name: hello-world-pod`のlabelを持ったpodと紐付けるという意味。

`port: 80`はこのserviceがうけつけるport。`targetPort: 8080`はpodが開放しているport。

lb-service.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-world-lb-svc
spec:
  type: LoadBalancer
  selector:
    name: hello-world-pod
  ports:
  - port: 80
    targetPort: 8080
```

serviceのデプロイ

```
$ kubectl create -f lb-service.yaml
```

`$ kubectl get services hello-world-lb-svc -w` のように `-w`を指定すると、変更や再登録があった際に更新された値が表示されて便利。(wはwatchのw)

EXTERNAL-IP が `<pending>` から有効な IP アドレス に変わったら(数分かかります)、そのアドレスにアクセスしてみましょう。Hello Worldが出力されるはずです。

```
$ curl http://$EXTERNAL_IP/
```

`type: LoadBalancer`を指定したので、load balancerも作成されたという次第です。

serviceには[他にも種類がある](https://qiita.com/sheepland/private/cdff472ba2d37784a125#service%E3%81%AE%E7%A8%AE%E9%A1%9E)。

ちなみに、以下のようにリソース名称は省略できる。

- `kubectl get pods`は`kubectl get po`
- `kubectl get deployments`は`kubectl get deploy`
- `kubectl get services`は`kubectl get svc`

省略名称は、`kubectl api-resources`と打つと`SHORTNAMES`の欄に表示される。
また`kubectl api-resources`で全リソースが表示される。


# 9, Ingressの作成

次はServiceとIngressを作って以下の構成にします。

![スクリーンショット 2018-10-14 22.13.15.png](https://qiita-image-store.s3.amazonaws.com/0/14124/77a66678-f351-d3e5-5db6-53888728d5b4.png)


Serviceはトランスポート層(L4)のプロキシを行うためのしくみです。
アプリケーション層(L7)のプロキシを行うためには、Ingressを使用する必要があります。
Ingressは以下が可能。

* load balancing
* SSLの終端
* URLのパスごとにどのserviceにリクエストを割り振るかの設定
* name-based virtual hosting (Apacheのバーチャルホスト相当)

service.yaml

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-world-svc
spec:
  type: NodePort
  selector:
    name: hello-world-pod
  ports:
    - port: 8080
```


ingress.yaml

```yaml
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: hello-world-ing
spec:
  rules:
  - http:
     paths:
      - path: /*
        backend:
          serviceName: hello-world-svc
          servicePort: 8080
```

`serviceName: hello-world-svc`で紐付けるserviceを指定しています。

Serviceの作成

```
$ kubectl create -f service.yaml
```

Ingressの作成

```
$ kubectl create -f ingress.yaml
```

確認

```
$ kubectl get ingress hello-world-ing -w
```


しばらくするとADDRESS欄に外部IPが付与される。

しかし外部IPが付与されたあとでもロードバランサがトラフィックを処理する準備が整うまで、HTTP 404 や HTTP 500 などのエラーが発生することがあります。正常にレスポンスが返るようになるまで結構かかる(10分以上?)ので気長に待ちましょう。。

`kubectl describe ingress`のレスポンス内のAnnotationsのbackendsのUnknownがHEALTHYになったら正常にレスポンスを返すようになる。

もしくはGCPのKubernetesの検出と負荷分散の画面でingressのステータスが正常になるまで待ちましょう

https://console.cloud.google.com/kubernetes/discovery

また、ロードバランサーは以下の画面から見れます

https://console.cloud.google.com/net-services/loadbalancing/loadBalancers/list

# 11, dashboard

dashboard機能もあります

以下を実行してdashboard機能をデプロイします

```
$ kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.0.0-beta1/aio/deploy/recommended.yaml
```

以下を実行して表示される文字列をコピーします

```
$ kubectl config view -o json | jq '.users[] | select(.name=="'$(kubectl config current-context)'") | .user."auth-provider".config."access-token"'
```

proxyを立ち上げます

```
$ kubectl proxy
```

そして、  http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/ へアクセス。

認証が必要になるので、以下の画面が表示されたら、「トークン」を選び、さっきコピった文字列を「Enter token」テキストボックスにペーストします。で、「サインイン」ボタンクリック。
<img width="874" alt="スクリーンショット 2019-07-13 16.42.32.png" src="https://qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/14124/30071fc0-e73f-13e7-f771-4958a032b6d9.png">



proxyコマンドを使う以外にも方法はあります。参考: https://qiita.com/sheepland/items/0ee17b80fcfb10227a41

# 11, 最後に

終わったらまとめて消しましょう。

```sh
$ kubectl delete service hello-world-lb-svc
$ kubectl delete service hello-world-svc
$ kubectl delete ingress hello-world-ing
$ kubectl delete deployment hello-world-deployment

# クラスタごと削除したい方はこちら。クラスタは残しておきたい人は後述するnodeサイズを0にする方法のほうがよい。
$ gcloud container clusters delete my-k8s
```

Kubernetesのリソースの全削除は以下でも可能

```
$ kubectl delete all --all
```

GCPのKubernetesはnodeの数だけGCEインスタンスを立ち上げています。なのでその分料金が発生しています。
しかしGCEの画面からインスタンスを削除してもKubernetesが自動的にインスタンスを復旧(立ち上げ直す)しちゃいます。
なので、以下のコマンドでnode数を0にするとよいです。こうするとclusterに登録したリソース(deploymentやserviceやingress等)の情報が失われません。num-nodesを再度例えば3にすると復旧します。

```
$ gcloud container clusters resize my-k8s --num-nodes=0
```


# その他のkubectlコマンド

https://kubernetes.io/docs/reference/kubectl/overview/#operations を参照

例えば以下がある

```
kubectl describe po
kubectl describe po POD名
kubectl logs pod名
kubectl explain pod
kubectl explain pod.spec
```

# Kubernetesをもっと知りたい、触りたい

* Kubernetes完全ガイド http://amzn.asia/d/9QFJw5X
    * これを買えば間違いがない
* 公式のWebでできるチュートリアル。Kubernetesクラスタを手元で立てずにチュートリアルができる。
    * https://kubernetes.io/docs/tutorials/kubernetes-basics/cluster-interactive/
* ローカルでKubernetesクラスタを簡単にたてるためのアプリケーション
    * https://github.com/kubernetes/minikube

