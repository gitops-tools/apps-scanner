package pipelines

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestParser(t *testing.T) {
	discoverTests := []struct {
		name  string
		items [][]corev1.Pod
		want  []Pipeline
		opts  []cmp.Option
	}{
		{
			name: "pods with no labels",
			items: [][]corev1.Pod{
				{
					makePod(),
				},
			},
			want: []Pipeline{},
		},
		{
			name: "one pipeline, one environment",
			items: [][]corev1.Pod{
				{
					makePod(withLabels(map[string]string{
						PipelineNameLabel:        "billing-pipeline",
						PipelineEnvironmentLabel: "test",
					})),
				},
			},
			want: []Pipeline{
				{
					Name:         "billing-pipeline",
					Environments: []string{"test"},
				},
			},
		},
		{
			name: "one pipeline, two environments",
			items: [][]corev1.Pod{
				{
					makePod(withLabels(map[string]string{
						PipelineNameLabel:        "billing-pipeline",
						PipelineEnvironmentLabel: "dev",
					})),
					makePod(withLabels(map[string]string{
						PipelineNameLabel:        "billing-pipeline",
						PipelineEnvironmentLabel: "staging",
					})),
				},
			},
			want: []Pipeline{
				{
					Name:         "billing-pipeline",
					Environments: []string{"dev", "staging"},
				},
			},
			opts: []cmp.Option{
				// This is needed in this case because there's no ordering
				// and so the test would be unstable.
				cmpopts.SortSlices(func(x, y string) bool {
					return strings.Compare(x, y) < 0
				}),
			},
		},
		{
			name: "one pipeline, two ordered environments",
			items: [][]corev1.Pod{
				{
					makePod(withLabels(map[string]string{
						PipelineNameLabel:        "billing-pipeline",
						PipelineEnvironmentLabel: "staging",
					})),
				},
				{
					makePod(withLabels(map[string]string{
						PipelineNameLabel:             "billing-pipeline",
						PipelineEnvironmentLabel:      "production",
						PipelineEnvironmentAfterLabel: "staging",
					})),
				},
			},
			want: []Pipeline{
				{
					Name:         "billing-pipeline",
					Environments: []string{"staging", "production"},
				},
			},
		},
	}

	for _, tt := range discoverTests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			for _, v := range tt.items {
				pods := &corev1.PodList{
					Items: v,
				}

				err := p.Add(pods)
				if err != nil {
					t.Fatal(err)
				}
			}

			pipelines, err := p.Pipelines()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, pipelines, tt.opts...); diff != "" {
				t.Fatalf("failed discovery:\n%s", diff)
			}
		})
	}
}

func TestParser_with_custom_labels(t *testing.T) {
	pods := []corev1.Pod{
		makePod(withLabels(map[string]string{
			"testing.pipeline":    "billing-pipeline",
			"testing.environment": "production",
			"testing.after":       "staging",
		})),
		makePod(withLabels(map[string]string{
			"testing.pipeline":    "billing-pipeline",
			"testing.environment": "staging",
		})),
	}

	p := NewParser(WithLabels("testing.pipeline", "testing.environment", "testing.after"))
	err := p.Add(&corev1.PodList{Items: pods})
	if err != nil {
		t.Fatal(err)
	}

	pipelines, err := p.Pipelines()
	if err != nil {
		t.Fatal(err)
	}

	want := []Pipeline{
		{
			Name:         "billing-pipeline",
			Environments: []string{"staging", "production"},
		},
	}

	if diff := cmp.Diff(want, pipelines); diff != "" {
		t.Fatalf("failed discovery:\n%s", diff)
	}
}

func makePod(opts ...func(runtime.Object)) corev1.Pod {
	p := corev1.Pod{}
	for _, o := range opts {
		o(&p)
	}
	return p
}

func withLabels(m map[string]string) func(runtime.Object) {
	accessor := meta.NewAccessor()
	return func(obj runtime.Object) {
		accessor.SetLabels(obj, m)
	}
}

func makeObjectMetaWithLabels(m map[string]string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Labels: m,
	}
}
