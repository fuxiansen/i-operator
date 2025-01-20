#下载kuberbuilder
1、version="v4.3.1"
   curl -L -o kubebuilder "https://github.com/kubernetes-sigs/kubebuilder/releases/download/${version}/kubebuilder_$(go env GOOS)_$(go env GOARCH)"
   chmod +x kubebuilder && mv kubebuilder /usr/local/bin/

#项目初始化
2、mkdir i-operator && cd i-operator
   kubebuilder init --domain crd.fuxiansen.com --repo github.com/fuxiansen/i-operator

#创建API
3、kubebuilder create api --group core --version v1 --kind Application --namespaced=true

#根据自定义的 CRD 生成对应的 yaml 文件
4、make manifests

#将 CRD 部署到集群
5、make install

#编译和运行控制器Controller
6、make run

#安装 CR 实例
7、kubectl apply -f config/samples/core_v1_application.yaml

# 构建并推送镜像，不包含 http:// 前缀
8、make docker-build docker-push IMG=192.168.126.106:81/fuxiansen/i-operator:latest

#根据 IMG 指定的镜像部署到集群
9、make deploy IMG=192.168.126.106:81/fuxiansen/i-operator:latest

#卸载 CRD
10、make uninstall

#删除控制器
11、make undeploy