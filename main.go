package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	_ "k8s.io/kubernetes/pkg/apis/apps/install"
	_ "k8s.io/kubernetes/pkg/apis/batch/install"
	_ "k8s.io/kubernetes/pkg/apis/core/install"
	_ "k8s.io/kubernetes/pkg/apis/extensions/install"

	//"github.com/appscode/go/log"
	_ "k8s.io/api/extensions/v1beta1"

	 //"k8s.io/client-go/kubernetes/scheme"
	//_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	"k8s.io/kubernetes/pkg/apis/apps"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"github.com/appscode/go/types"
	oneliner "github.com/the-redback/go-oneliners"
	appsv1 "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	TEST_HEADLESS_SERVICE    = "headless"
	TestSourceDataVolumeName = "source-data"
	TestSourceDataMountPath  = "/source/data"
	v1Version                = "v1"
	v1beta1Version           = "v1beta1"
	ReplicationController    = "ReplicationController"
	StatefulSet	= "StatefulSet"
)

func main() {
	v1Obj := []byte("saldfjl")
	sgvk := schema.GroupVersionKind{Group: appsv1.GroupName, Version: v1Version, Kind: StatefulSet}
	fmt.Println("Group:", sgvk.Group, " Version:", sgvk.Version, " Kind:", sgvk.Kind)
	v1beta1Obj:=&appsv1beta1.StatefulSet{}
	err := Convert(sgvk,  v1Obj,v1beta1Obj)
	if err!=nil{
		log.Panic(err)
	}
	oneliner.PrettyJson(v1beta1Obj, "V1beta1Object")
}

func Convert(sgvk schema.GroupVersionKind, in,out interface{}) error {

	var sourceObj, internalObj interface{}
	switch sgvk.Group {
	case apps.GroupName:
		switch sgvk.Version {
		case v1Version:
			switch sgvk.Kind {
			case StatefulSet:
				sourceObj=&appsv1.StatefulSet{}
			default:
				return fmt.Errorf("Unknown Kind")
			}
		case v1beta1Version:


		}
		internalObj=&apps.StatefulSet{}
	case core.GroupName:
	case extensions.GroupName:
	default:
		return fmt.Errorf("Unkown Group")
	}
	rt:=reflect.TypeOf(in)
	if rt.Kind() == reflect.Ptr{
		sourceObj = in
	} else {
		json.Unmarshal(in.([]byte),sourceObj)
	}
	oneliner.PrettyJson(sourceObj,"Source Object")
	err:=legacyscheme.Scheme.Convert(sourceObj,internalObj,nil)
	if err!=nil{
		return err
	}

	oneliner.PrettyJson(internalObj,"InternalObj")
	err = legacyscheme.Scheme.Convert(sourceObj,out,nil)
	if err!=nil{
		return err
	}
	//internalObj := &apps.StatefulSet{}
	//err := legacyscheme.Scheme.Convert(v1Obj, internalObj, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//oneliner.PrettyJson(internalObj,"internalObj")
	//
	//
	//err = legacyscheme.Scheme.Convert(internalObj, v1beta1Obj, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//oneliner.PrettyJson(v1beta1Obj,"v1beta1Obj")
	return nil
}

func getStatefulSet() *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "stash-demo",
			Labels: map[string]string{
				"app": "stash-demo",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: types.Int32P(1),
			Template: PodTemplate(),
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: appsv1beta1.RollingUpdateStatefulSetStrategyType,
			},
		},
	}
}

func getReplicationController() corev1.ReplicationController {
	podTemplate := PodTemplate()
	return corev1.ReplicationController{
		ObjectMeta: metav1.ObjectMeta{
			Name: "stash-demo",
			Labels: map[string]string{
				"app": "stash-demo",
			},
		},
		Spec: corev1.ReplicationControllerSpec{
			Replicas: types.Int32P(1),
			Template: &podTemplate,
		},
	}
}

func ReplicaSet() extensionsv1beta1.ReplicaSet {
	return extensionsv1beta1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: "stash-demo",
			Labels: map[string]string{
				"app": "stash-demo",
			},
		},
		Spec: extensionsv1beta1.ReplicaSetSpec{
			Replicas: types.Int32P(1),
			Template: PodTemplate(),
		},
	}
}

func PodTemplate() corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "stash-demo",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            "busybox",
					Image:           "busybox",
					ImagePullPolicy: corev1.PullIfNotPresent,
					Command: []string{
						"sleep",
						"3600",
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      TestSourceDataVolumeName,
							MountPath: TestSourceDataMountPath,
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: TestSourceDataVolumeName,
					VolumeSource: corev1.VolumeSource{
						GitRepo: &corev1.GitRepoVolumeSource{
							Repository: "https://github.com/appscode/stash-data.git",
						},
					},
				},
			},
		},
	}
}
