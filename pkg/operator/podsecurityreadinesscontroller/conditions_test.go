package podsecurityreadinesscontroller

import (
	"testing"

	operatorv1 "github.com/openshift/api/operator/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCondition(t *testing.T) {
	t.Run("with namespaces", func(t *testing.T) {
		namespaces := []string{"namespace1", "namespace2"}
		expectedCondition := operatorv1.OperatorCondition{
			Type:    PodSecurityCustomerType,
			Status:  operatorv1.ConditionTrue,
			Reason:  "PSViolationsDetected",
			Message: "Violations detected in namespaces: [namespace1 namespace2]",
		}

		condition := makeCondition(PodSecurityCustomerType, namespaces)

		if condition.Type != expectedCondition.Type {
			t.Errorf("expected condition type %s, got %s", expectedCondition.Type, condition.Type)
		}

		if condition.Status != expectedCondition.Status {
			t.Errorf("expected condition status %s, got %s", expectedCondition.Status, condition.Status)
		}

		if condition.Reason != expectedCondition.Reason {
			t.Errorf("expected condition reason %s, got %s", expectedCondition.Reason, condition.Reason)
		}

		if condition.Message != expectedCondition.Message {
			t.Errorf("expected condition message %s, got %s", expectedCondition.Message, condition.Message)
		}
	})

	t.Run("without namespaces", func(t *testing.T) {
		namespaces := []string{}
		expectedCondition := operatorv1.OperatorCondition{
			Type:   PodSecurityCustomerType,
			Status: operatorv1.ConditionFalse,
			Reason: "ExpectedReason",
		}

		condition := makeCondition(PodSecurityCustomerType, namespaces)

		if condition.Type != expectedCondition.Type {
			t.Errorf("expected condition type %s, got %s", expectedCondition.Type, condition.Type)
		}

		if condition.Status != expectedCondition.Status {
			t.Errorf("expected condition status %s, got %s", expectedCondition.Status, condition.Status)
		}

		if condition.Reason != expectedCondition.Reason {
			t.Errorf("expected condition reason %s, got %s", expectedCondition.Reason, condition.Reason)
		}

		if condition.Message != expectedCondition.Message {
			t.Errorf("expected condition message %s, got %s", expectedCondition.Message, condition.Message)
		}
	})

}

func TestOperatorStatus(t *testing.T) {
	for _, tt := range []struct {
		name      string
		namespace []*corev1.Namespace
		expected  map[string]operatorv1.ConditionStatus
	}{
		{
			name: "with default namespace",
			namespace: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "syncer-by-default",
					},
				},
			},
			expected: map[string]operatorv1.ConditionStatus{
				"PodSecurityCustomerEvaluationConditionsDetected":       operatorv1.ConditionTrue,
				"PodSecurityOpenshiftEvaluationConditionsDetected":      operatorv1.ConditionFalse,
				"PodSecurityRunLevelZeroEvaluationConditionsDetected":   operatorv1.ConditionFalse,
				"PodSecurityDisabledSyncerEvaluationConditionsDetected": operatorv1.ConditionFalse,
			},
		},
		{
			name: "with customer disabled syncer",
			namespace: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "syncer-no-thx",
						Labels: map[string]string{
							"security.openshift.io/scc.podSecurityLabelSync": "false",
						},
					},
				},
			},
			expected: map[string]operatorv1.ConditionStatus{
				"PodSecurityCustomerEvaluationConditionsDetected":       operatorv1.ConditionFalse,
				"PodSecurityOpenshiftEvaluationConditionsDetected":      operatorv1.ConditionFalse,
				"PodSecurityRunLevelZeroEvaluationConditionsDetected":   operatorv1.ConditionFalse,
				"PodSecurityDisabledSyncerEvaluationConditionsDetected": operatorv1.ConditionTrue,
			},
		},
		{
			name: "with customer re-enabled syncer",
			namespace: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "syncer-yes-plz",
						Labels: map[string]string{
							"security.openshift.io/scc.podSecurityLabelSync": "true",
						},
					},
				},
			},
			expected: map[string]operatorv1.ConditionStatus{
				"PodSecurityCustomerEvaluationConditionsDetected":       operatorv1.ConditionTrue,
				"PodSecurityOpenshiftEvaluationConditionsDetected":      operatorv1.ConditionFalse,
				"PodSecurityRunLevelZeroEvaluationConditionsDetected":   operatorv1.ConditionFalse,
				"PodSecurityDisabledSyncerEvaluationConditionsDetected": operatorv1.ConditionFalse,
			},
		},
		{
			name: "with openshift namespace",
			namespace: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "openshift-fail",
					},
				},
			},
			expected: map[string]operatorv1.ConditionStatus{
				"PodSecurityCustomerEvaluationConditionsDetected":       operatorv1.ConditionFalse,
				"PodSecurityOpenshiftEvaluationConditionsDetected":      operatorv1.ConditionTrue,
				"PodSecurityRunLevelZeroEvaluationConditionsDetected":   operatorv1.ConditionFalse,
				"PodSecurityDisabledSyncerEvaluationConditionsDetected": operatorv1.ConditionFalse,
			},
		},
		{
			name: "with run-level 0 namespace",
			namespace: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "kube-system",
					},
				},
			},
			expected: map[string]operatorv1.ConditionStatus{
				"PodSecurityCustomerEvaluationConditionsDetected":       operatorv1.ConditionFalse,
				"PodSecurityOpenshiftEvaluationConditionsDetected":      operatorv1.ConditionFalse,
				"PodSecurityRunLevelZeroEvaluationConditionsDetected":   operatorv1.ConditionTrue,
				"PodSecurityDisabledSyncerEvaluationConditionsDetected": operatorv1.ConditionFalse,
			},
		},
		{
			name: "with other customer types in combination",
			namespace: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "foobar",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "foobar",
						Labels: map[string]string{
							"security.openshift.io/scc.podSecurityLabelSync": "false",
						},
					},
				},
			},
			expected: map[string]operatorv1.ConditionStatus{
				"PodSecurityCustomerEvaluationConditionsDetected":       operatorv1.ConditionTrue,
				"PodSecurityOpenshiftEvaluationConditionsDetected":      operatorv1.ConditionFalse,
				"PodSecurityRunLevelZeroEvaluationConditionsDetected":   operatorv1.ConditionFalse,
				"PodSecurityDisabledSyncerEvaluationConditionsDetected": operatorv1.ConditionTrue,
			},
		},
		{
			name: "with other system types in combination",
			namespace: []*corev1.Namespace{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "openshift-namespace",
						Labels: map[string]string{
							"pod-security.kubernetes.io/audit":         "restricted",
							"pod-security.kubernetes.io/audit-version": "v1.24",
							"pod-security.kubernetes.io/warn":          "restricted",
							"pod-security.kubernetes.io/warn-version":  "v1.24",
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:   "kube-system",
						Labels: map[string]string{},
					},
				},
			},
			expected: map[string]operatorv1.ConditionStatus{
				"PodSecurityCustomerEvaluationConditionsDetected":       operatorv1.ConditionFalse,
				"PodSecurityOpenshiftEvaluationConditionsDetected":      operatorv1.ConditionTrue,
				"PodSecurityRunLevelZeroEvaluationConditionsDetected":   operatorv1.ConditionTrue,
				"PodSecurityDisabledSyncerEvaluationConditionsDetected": operatorv1.ConditionFalse,
			},
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			cond := podSecurityOperatorConditions{}

			for _, ns := range tt.namespace {
				cond.addViolation(ns)
			}

			status := &operatorv1.OperatorStatus{}
			for _, f := range cond.toConditionFuncs() {
				if err := f(status); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			for expectedType, expectedStatus := range tt.expected {
				found := false

				for _, condition := range status.Conditions {
					if condition.Type == expectedType {
						found = true
						if condition.Status != expectedStatus {
							t.Errorf("expected %s to be %v, have %v", expectedType, expectedStatus, condition.Status)
						}
					}
				}

				if !found {
					t.Errorf("expected condition %s not found", expectedType)
				}
			}
		})
	}
}
