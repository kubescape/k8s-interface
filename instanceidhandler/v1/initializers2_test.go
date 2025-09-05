package instanceidhandler

import (
	"fmt"
	"testing"

	"github.com/kubescape/k8s-interface/instanceidhandler"
	"github.com/kubescape/k8s-interface/instanceidhandler/v1/containerinstance"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestGenerateInstanceIDFromRuntime2(t *testing.T) {
	type args struct {
		w         runtime.Object
		jsonPaths []string
	}
	tests := []struct {
		name    string
		args    args
		want    []instanceidhandler.IInstanceID
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "DaemonSet",
			args: args{
				w: &appsv1.DaemonSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "daemonset1",
						Namespace: "default",
					},
					Spec: appsv1.DaemonSetSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "container1",
									},
								},
							},
						},
					},
				},
			},
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "DaemonSet",
					Name:          "daemonset1",
					ContainerName: "container1",
					InstanceType:  "container",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Deployment",
			args: args{
				w: &appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "deployment1",
						Namespace: "default",
					},
					Spec: appsv1.DeploymentSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "container1",
									},
								},
							},
						},
					},
				},
			},
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "Deployment",
					Name:          "deployment1",
					ContainerName: "container1",
					InstanceType:  "container",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ReplicaSet",
			args: args{
				w: &appsv1.ReplicaSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "replicaset1",
						Namespace: "default",
					},
					Spec: appsv1.ReplicaSetSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "container1",
									},
								},
							},
						},
					},
				},
			},
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "ReplicaSet",
					Name:          "replicaset1",
					ContainerName: "container1",
					InstanceType:  "container",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "StatefulSet",
			args: args{
				w: &appsv1.StatefulSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "statefulset1",
						Namespace: "default",
					},
					Spec: appsv1.StatefulSetSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "container1",
									},
								},
							},
						},
					},
				},
			},
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "apps/v1",
					Namespace:     "default",
					Kind:          "StatefulSet",
					Name:          "statefulset1",
					ContainerName: "container1",
					InstanceType:  "container",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "CronJob",
			args: args{
				w: &batchv1.CronJob{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "cronjob1",
						Namespace: "default",
					},
					Spec: batchv1.CronJobSpec{
						JobTemplate: batchv1.JobTemplateSpec{
							Spec: batchv1.JobSpec{
								Template: v1.PodTemplateSpec{
									Spec: v1.PodSpec{
										Containers: []v1.Container{
											{
												Name: "container1",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "batch/v1",
					Namespace:     "default",
					Kind:          "CronJob",
					Name:          "cronjob1",
					ContainerName: "container1",
					InstanceType:  "container",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Job",
			args: args{
				w: &batchv1.Job{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "job1",
						Namespace: "default",
					},
					Spec: batchv1.JobSpec{
						Template: v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "container1",
									},
								},
							},
						},
					},
				},
			},
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "batch/v1",
					Namespace:     "default",
					Kind:          "Job",
					Name:          "job1",
					ContainerName: "container1",
					InstanceType:  "container",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Pod",
			args: args{
				w: &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "pod1",
						Namespace: "default",
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Name: "container1",
							},
						},
					},
				},
			},
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "v1",
					Namespace:     "default",
					Kind:          "Pod",
					Name:          "pod1",
					ContainerName: "container1",
					InstanceType:  "container",
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "PodTemplate",
			args: args{
				w: &v1.PodTemplate{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "podtemplate1",
						Namespace: "default",
					},
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Containers: []v1.Container{
								{
									Name: "container1",
								},
							},
						},
					},
				},
			},
			want: []instanceidhandler.IInstanceID{
				&containerinstance.InstanceID{
					ApiVersion:    "v1",
					Namespace:     "default",
					Kind:          "PodTemplate",
					Name:          "podtemplate1",
					ContainerName: "container1",
					InstanceType:  "container",
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateInstanceIDFromRuntimeObj(tt.args.w, tt.args.jsonPaths)
			if !tt.wantErr(t, err, fmt.Sprintf("GenerateInstanceIDFromRuntimeObj(%v, %v)", tt.args.w, tt.args.jsonPaths)) {
				return
			}
			assert.Equalf(t, tt.want, got, "GenerateInstanceIDFromRuntimeObj(%v, %v)", tt.args.w, tt.args.jsonPaths)
		})
	}
}
