package cnv

import (
	"bytes"
	"html/template"
)

// Manifests returns manifests needed to deploy CNV
func Manifests(OpenShiftVersion string) (map[string]string, error) {
	manifests := make(map[string]string)
	manifests["99_openshift-cnv_crd.yaml"] = cnvCRD
	manifests["99_openshift-cnv_ns.yaml"] = cnvNamespace
	cnvSubs, err := subscription(OpenShiftVersion)
	if err != nil {
		return map[string]string{}, err
	}
	manifests["99_openshift-cnv_subscription.yaml"] = cnvSubs
	manifests["99_openshift-cnv_operator_group.yaml"] = cnvGroup
	manifests["99_openshift-cnv_hco.yaml"] = cnvHCO
	return manifests, nil
}

func subscription(OpenShiftVersion string) (string, error) {
	data := map[string]string{
		"OPENSHIFT_VERSION": OpenShiftVersion,
	}
	tmpl, err := template.New("cnvSubscription").Parse(cnvSubscription)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

const cnvSubscription = `apiVersion: operators.coreos.com/v1alpha1
kind: Subscription
metadata:
  name: hco-operatorhub
  namespace: openshift-cnv
spec:
  source: redhat-operators
  sourceNamespace: openshift-marketplace
  name: kubevirt-hyperconverged
  channel: "{{.OPENSHIFT_VERSION}}"
  installPlanApproval: "Automatic"`

const cnvNamespace = `apiVersion: v1
kind: Namespace
metadata:
  name: openshift-cnv`

const cnvGroup = `apiVersion: operators.coreos.com/v1
kind: OperatorGroup
metadata:
  name: kubevirt-hyperconverged-group
  namespace: openshift-cnv
spec:
  targetNamespaces:
	- openshift-cnv`

const cnvHCO = `apiVersion: hco.kubevirt.io/v1beta1
kind: HyperConverged
metadata:
  name: kubevirt-hyperconverged
  namespace: openshift-cnv
spec:
  BareMetalPlatform: true`

const cnvCRD = `apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  labels:
    operators.coreos.com/kubevirt-hyperconverged.openshift-cnv: ""
  name: hyperconvergeds.hco.kubevirt.io
spec:
  conversion:
    strategy: None
  group: hco.kubevirt.io
  names:
    categories:
    - all
    kind: HyperConverged
    listKind: HyperConvergedList
    plural: hyperconvergeds
    shortNames:
    - hco
    - hcos
    singular: hyperconverged
  scope: Namespaced
  versions:
  - additionalPrinterColumns:
    - jsonPath: .metadata.creationTimestamp
      name: Age
      type: date
    name: v1alpha1
    schema:
      openAPIV3Schema:
        description: HyperConverged is the Schema for the hyperconvergeds API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            properties:
              name:
                pattern: kubevirt-hyperconverged
                type: string
            type: object
          spec:
            description: HyperConvergedSpec defines the desired state of HyperConverged
            properties:
              BareMetalPlatform:
                description: BareMetalPlatform indicates whether the infrastructure
                  is baremetal.
                type: boolean
              LocalStorageClassName:
                description: LocalStorageClassName the name of the local storage class.
                type: string
              version:
                description: operator version
                type: string
            type: object
          status:
            description: HyperConvergedStatus defines the observed state of HyperConverged
            properties:
              conditions:
                description: Conditions describes the state of the HyperConverged
                  resource.
                items:
                  description: Condition represents the state of the operator's reconciliation
                    functionality.
                  properties:
                    lastHeartbeatTime:
                      format: date-time
                      type: string
                    lastTransitionTime:
                      format: date-time
                      type: string
                    message:
                      type: string
                    reason:
                      type: string
                    status:
                      type: string
                    type:
                      description: ConditionType is the state of the operator's reconciliation
                        functionality.
                      type: string
                  required:
                  - status
                  - type
                  type: object
                type: array
              relatedObjects:
                description: RelatedObjects is a list of objects created and maintained
                  by this operator. Object references will be added to this list after
                  they have been created AND found in the cluster.
                items:
                  description: 'ObjectReference contains enough information to let
                    you inspect or modify the referred object. --- New uses of this
                    type are discouraged because of difficulty describing its usage
                    when embedded in APIs.  1. Ignored fields.  It includes many fields
                    which are not generally honored.  For instance, ResourceVersion
                    and FieldPath are both very rarely valid in actual usage.  2.
                    Invalid usage help.  It is impossible to add specific help for
                    individual usage.  In most embedded usages, there are particular     restrictions
                    like, "must refer only to types A and B" or "UID not honored"
                    or "name must be restricted".     Those cannot be well described
                    when embedded.  3. Inconsistent validation.  Because the usages
                    are different, the validation rules are different by usage, which
                    makes it hard for users to predict what will happen.  4. The fields
                    are both imprecise and overly precise.  Kind is not a precise
                    mapping to a URL. This can produce ambiguity     during interpretation
                    and require a REST mapping.  In most cases, the dependency is
                    on the group,resource tuple     and the version of the actual
                    struct is irrelevant.  5. We cannot easily change it.  Because
                    this type is embedded in many locations, updates to this type     will
                    affect numerous schemas.  Don''t make new APIs embed an underspecified
                    API type they do not control. Instead of using this type, create
                    a locally provided and used type that is well-focused on your
                    reference. For example, ServiceReferences for admission registration:
                    https://github.com/kubernetes/api/blob/release-1.17/admissionregistration/v1/types.go#L533
                    .'
                  properties:
                    apiVersion:
                      description: API version of the referent.
                      type: string
                    fieldPath:
                      description: 'If referring to a piece of an object instead of
                        an entire object, this string should contain a valid JSON/Go
                        field access statement, such as desiredState.manifest.containers[2].
                        For example, if the object reference is to a container within
                        a pod, this would take on a value like: "spec.containers{name}"
                        (where "name" refers to the name of the container that triggered
                        the event) or if no container name is specified "spec.containers[2]"
                        (container with index 2 in this pod). This syntax is chosen
                        only to have some well-defined way of referencing a part of
                        an object. TODO: this design is not final and this field is
                        subject to change in the future.'
                      type: string
                    kind:
                      description: 'Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                      type: string
                    name:
                      description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names'
                      type: string
                    namespace:
                      description: 'Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/'
                      type: string
                    resourceVersion:
                      description: 'Specific resourceVersion to which this reference
                        is made, if any. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency'
                      type: string
                    uid:
                      description: 'UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids'
                      type: string
                  type: object
                type: array
              versions:
                description: 'Versions is a list of HCO component versions, as name/version
                  pairs. The version with a name of "operator" is the HCO version
                  itself, as described here: https://github.com/openshift/cluster-version-operator/blob/master/docs/dev/clusteroperator.md#version'
                items:
                  properties:
                    name:
                      type: string
                    version:
                      type: string
                  type: object
                type: array
            type: object
        type: object
    served: true
    storage: false
    subresources:
      status: {}
  - additionalPrinterColumns:
    - jsonPath: .metadata.creationTimestamp
      name: Age
      type: date
    name: v1beta1
    schema:
      openAPIV3Schema:
        description: HyperConverged is the Schema for the hyperconvergeds API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            properties:
              name:
                pattern: kubevirt-hyperconverged
                type: string
            type: object
          spec:
            description: HyperConvergedSpec defines the desired state of HyperConverged
            properties:
              bareMetalPlatform:
                description: BareMetalPlatform indicates whether the infrastructure
                  is baremetal.
                type: boolean
              featureGates:
                description: featureGates is a map of feature gate flags. Setting
                  a flag to 'true' will enable the feature. Setting 'false' or removing
                  the feature gate, disables the feature.
                properties:
                  hotplugVolumes:
                    description: Allow attaching a data volume to a running VMI
                    type: boolean
                type: object
              infra:
                description: infra HyperConvergedConfig influences the pod configuration
                  (currently only placement) for all the infra components needed on
                  the virtualization enabled cluster but not necessarely directly
                  on each node running VMs/VMIs.
                properties:
                  nodePlacement:
                    description: NodePlacement describes node scheduling configuration.
                    properties:
                      affinity:
                        description: affinity enables pod affinity/anti-affinity placement
                          expanding the types of constraints that can be expressed
                          with nodeSelector. affinity is going to be applied to the
                          relevant kind of pods in parallel with nodeSelector See
                          https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
                        properties:
                          nodeAffinity:
                            description: Describes node affinity scheduling rules
                              for the pod.
                            properties:
                              preferredDuringSchedulingIgnoredDuringExecution:
                                description: The scheduler will prefer to schedule
                                  pods to nodes that satisfy the affinity expressions
                                  specified by this field, but it may choose a node
                                  that violates one or more of the expressions. The
                                  node that is most preferred is the one with the
                                  greatest sum of weights, i.e. for each node that
                                  meets all of the scheduling requirements (resource
                                  request, requiredDuringScheduling affinity expressions,
                                  etc.), compute a sum by iterating through the elements
                                  of this field and adding "weight" to the sum if
                                  the node matches the corresponding matchExpressions;
                                  the node(s) with the highest sum are the most preferred.
                                items:
                                  description: An empty preferred scheduling term
                                    matches all objects with implicit weight 0 (i.e.
                                    it's a no-op). A null preferred scheduling term
                                    matches no objects (i.e. is also a no-op).
                                  properties:
                                    preference:
                                      description: A node selector term, associated
                                        with the corresponding weight.
                                      properties:
                                        matchExpressions:
                                          description: A list of node selector requirements
                                            by node's labels.
                                          items:
                                            description: A node selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: The label key that the
                                                  selector applies to.
                                                type: string
                                              operator:
                                                description: Represents a key's relationship
                                                  to a set of values. Valid operators
                                                  are In, NotIn, Exists, DoesNotExist.
                                                  Gt, and Lt.
                                                type: string
                                              values:
                                                description: An array of string values.
                                                  If the operator is In or NotIn,
                                                  the values array must be non-empty.
                                                  If the operator is Exists or DoesNotExist,
                                                  the values array must be empty.
                                                  If the operator is Gt or Lt, the
                                                  values array must have a single
                                                  element, which will be interpreted
                                                  as an integer. This array is replaced
                                                  during a strategic merge patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                        matchFields:
                                          description: A list of node selector requirements
                                            by node's fields.
                                          items:
                                            description: A node selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: The label key that the
                                                  selector applies to.
                                                type: string
                                              operator:
                                                description: Represents a key's relationship
                                                  to a set of values. Valid operators
                                                  are In, NotIn, Exists, DoesNotExist.
                                                  Gt, and Lt.
                                                type: string
                                              values:
                                                description: An array of string values.
                                                  If the operator is In or NotIn,
                                                  the values array must be non-empty.
                                                  If the operator is Exists or DoesNotExist,
                                                  the values array must be empty.
                                                  If the operator is Gt or Lt, the
                                                  values array must have a single
                                                  element, which will be interpreted
                                                  as an integer. This array is replaced
                                                  during a strategic merge patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                      type: object
                                    weight:
                                      description: Weight associated with matching
                                        the corresponding nodeSelectorTerm, in the
                                        range 1-100.
                                      format: int32
                                      type: integer
                                  required:
                                  - preference
                                  - weight
                                  type: object
                                type: array
                              requiredDuringSchedulingIgnoredDuringExecution:
                                description: If the affinity requirements specified
                                  by this field are not met at scheduling time, the
                                  pod will not be scheduled onto the node. If the
                                  affinity requirements specified by this field cease
                                  to be met at some point during pod execution (e.g.
                                  due to an update), the system may or may not try
                                  to eventually evict the pod from its node.
                                properties:
                                  nodeSelectorTerms:
                                    description: Required. A list of node selector
                                      terms. The terms are ORed.
                                    items:
                                      description: A null or empty node selector term
                                        matches no objects. The requirements of them
                                        are ANDed. The TopologySelectorTerm type implements
                                        a subset of the NodeSelectorTerm.
                                      properties:
                                        matchExpressions:
                                          description: A list of node selector requirements
                                            by node's labels.
                                          items:
                                            description: A node selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: The label key that the
                                                  selector applies to.
                                                type: string
                                              operator:
                                                description: Represents a key's relationship
                                                  to a set of values. Valid operators
                                                  are In, NotIn, Exists, DoesNotExist.
                                                  Gt, and Lt.
                                                type: string
                                              values:
                                                description: An array of string values.
                                                  If the operator is In or NotIn,
                                                  the values array must be non-empty.
                                                  If the operator is Exists or DoesNotExist,
                                                  the values array must be empty.
                                                  If the operator is Gt or Lt, the
                                                  values array must have a single
                                                  element, which will be interpreted
                                                  as an integer. This array is replaced
                                                  during a strategic merge patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                        matchFields:
                                          description: A list of node selector requirements
                                            by node's fields.
                                          items:
                                            description: A node selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: The label key that the
                                                  selector applies to.
                                                type: string
                                              operator:
                                                description: Represents a key's relationship
                                                  to a set of values. Valid operators
                                                  are In, NotIn, Exists, DoesNotExist.
                                                  Gt, and Lt.
                                                type: string
                                              values:
                                                description: An array of string values.
                                                  If the operator is In or NotIn,
                                                  the values array must be non-empty.
                                                  If the operator is Exists or DoesNotExist,
                                                  the values array must be empty.
                                                  If the operator is Gt or Lt, the
                                                  values array must have a single
                                                  element, which will be interpreted
                                                  as an integer. This array is replaced
                                                  during a strategic merge patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                      type: object
                                    type: array
                                required:
                                - nodeSelectorTerms
                                type: object
                            type: object
                          podAffinity:
                            description: Describes pod affinity scheduling rules (e.g.
                              co-locate this pod in the same node, zone, etc. as some
                              other pod(s)).
                            properties:
                              preferredDuringSchedulingIgnoredDuringExecution:
                                description: The scheduler will prefer to schedule
                                  pods to nodes that satisfy the affinity expressions
                                  specified by this field, but it may choose a node
                                  that violates one or more of the expressions. The
                                  node that is most preferred is the one with the
                                  greatest sum of weights, i.e. for each node that
                                  meets all of the scheduling requirements (resource
                                  request, requiredDuringScheduling affinity expressions,
                                  etc.), compute a sum by iterating through the elements
                                  of this field and adding "weight" to the sum if
                                  the node has pods which matches the corresponding
                                  podAffinityTerm; the node(s) with the highest sum
                                  are the most preferred.
                                items:
                                  description: The weights of all of the matched WeightedPodAffinityTerm
                                    fields are added per-node to find the most preferred
                                    node(s)
                                  properties:
                                    podAffinityTerm:
                                      description: Required. A pod affinity term,
                                        associated with the corresponding weight.
                                      properties:
                                        labelSelector:
                                          description: A label query over a set of
                                            resources, in this case pods.
                                          properties:
                                            matchExpressions:
                                              description: matchExpressions is a list
                                                of label selector requirements. The
                                                requirements are ANDed.
                                              items:
                                                description: A label selector requirement
                                                  is a selector that contains values,
                                                  a key, and an operator that relates
                                                  the key and values.
                                                properties:
                                                  key:
                                                    description: key is the label
                                                      key that the selector applies
                                                      to.
                                                    type: string
                                                  operator:
                                                    description: operator represents
                                                      a key's relationship to a set
                                                      of values. Valid operators are
                                                      In, NotIn, Exists and DoesNotExist.
                                                    type: string
                                                  values:
                                                    description: values is an array
                                                      of string values. If the operator
                                                      is In or NotIn, the values array
                                                      must be non-empty. If the operator
                                                      is Exists or DoesNotExist, the
                                                      values array must be empty.
                                                      This array is replaced during
                                                      a strategic merge patch.
                                                    items:
                                                      type: string
                                                    type: array
                                                required:
                                                - key
                                                - operator
                                                type: object
                                              type: array
                                            matchLabels:
                                              additionalProperties:
                                                type: string
                                              description: matchLabels is a map of
                                                {key,value} pairs. A single {key,value}
                                                in the matchLabels map is equivalent
                                                to an element of matchExpressions,
                                                whose key field is "key", the operator
                                                is "In", and the values array contains
                                                only "value". The requirements are
                                                ANDed.
                                              type: object
                                          type: object
                                        namespaces:
                                          description: namespaces specifies which
                                            namespaces the labelSelector applies to
                                            (matches against); null or empty list
                                            means "this pod's namespace"
                                          items:
                                            type: string
                                          type: array
                                        topologyKey:
                                          description: This pod should be co-located
                                            (affinity) or not co-located (anti-affinity)
                                            with the pods matching the labelSelector
                                            in the specified namespaces, where co-located
                                            is defined as running on a node whose
                                            value of the label with key topologyKey
                                            matches that of any node on which any
                                            of the selected pods is running. Empty
                                            topologyKey is not allowed.
                                          type: string
                                      required:
                                      - topologyKey
                                      type: object
                                    weight:
                                      description: weight associated with matching
                                        the corresponding podAffinityTerm, in the
                                        range 1-100.
                                      format: int32
                                      type: integer
                                  required:
                                  - podAffinityTerm
                                  - weight
                                  type: object
                                type: array
                              requiredDuringSchedulingIgnoredDuringExecution:
                                description: If the affinity requirements specified
                                  by this field are not met at scheduling time, the
                                  pod will not be scheduled onto the node. If the
                                  affinity requirements specified by this field cease
                                  to be met at some point during pod execution (e.g.
                                  due to a pod label update), the system may or may
                                  not try to eventually evict the pod from its node.
                                  When there are multiple elements, the lists of nodes
                                  corresponding to each podAffinityTerm are intersected,
                                  i.e. all terms must be satisfied.
                                items:
                                  description: Defines a set of pods (namely those
                                    matching the labelSelector relative to the given
                                    namespace(s)) that this pod should be co-located
                                    (affinity) or not co-located (anti-affinity) with,
                                    where co-located is defined as running on a node
                                    whose value of the label with key <topologyKey>
                                    matches that of any node on which a pod of the
                                    set of pods is running
                                  properties:
                                    labelSelector:
                                      description: A label query over a set of resources,
                                        in this case pods.
                                      properties:
                                        matchExpressions:
                                          description: matchExpressions is a list
                                            of label selector requirements. The requirements
                                            are ANDed.
                                          items:
                                            description: A label selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: key is the label key
                                                  that the selector applies to.
                                                type: string
                                              operator:
                                                description: operator represents a
                                                  key's relationship to a set of values.
                                                  Valid operators are In, NotIn, Exists
                                                  and DoesNotExist.
                                                type: string
                                              values:
                                                description: values is an array of
                                                  string values. If the operator is
                                                  In or NotIn, the values array must
                                                  be non-empty. If the operator is
                                                  Exists or DoesNotExist, the values
                                                  array must be empty. This array
                                                  is replaced during a strategic merge
                                                  patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                        matchLabels:
                                          additionalProperties:
                                            type: string
                                          description: matchLabels is a map of {key,value}
                                            pairs. A single {key,value} in the matchLabels
                                            map is equivalent to an element of matchExpressions,
                                            whose key field is "key", the operator
                                            is "In", and the values array contains
                                            only "value". The requirements are ANDed.
                                          type: object
                                      type: object
                                    namespaces:
                                      description: namespaces specifies which namespaces
                                        the labelSelector applies to (matches against);
                                        null or empty list means "this pod's namespace"
                                      items:
                                        type: string
                                      type: array
                                    topologyKey:
                                      description: This pod should be co-located (affinity)
                                        or not co-located (anti-affinity) with the
                                        pods matching the labelSelector in the specified
                                        namespaces, where co-located is defined as
                                        running on a node whose value of the label
                                        with key topologyKey matches that of any node
                                        on which any of the selected pods is running.
                                        Empty topologyKey is not allowed.
                                      type: string
                                  required:
                                  - topologyKey
                                  type: object
                                type: array
                            type: object
                          podAntiAffinity:
                            description: Describes pod anti-affinity scheduling rules
                              (e.g. avoid putting this pod in the same node, zone,
                              etc. as some other pod(s)).
                            properties:
                              preferredDuringSchedulingIgnoredDuringExecution:
                                description: The scheduler will prefer to schedule
                                  pods to nodes that satisfy the anti-affinity expressions
                                  specified by this field, but it may choose a node
                                  that violates one or more of the expressions. The
                                  node that is most preferred is the one with the
                                  greatest sum of weights, i.e. for each node that
                                  meets all of the scheduling requirements (resource
                                  request, requiredDuringScheduling anti-affinity
                                  expressions, etc.), compute a sum by iterating through
                                  the elements of this field and adding "weight" to
                                  the sum if the node has pods which matches the corresponding
                                  podAffinityTerm; the node(s) with the highest sum
                                  are the most preferred.
                                items:
                                  description: The weights of all of the matched WeightedPodAffinityTerm
                                    fields are added per-node to find the most preferred
                                    node(s)
                                  properties:
                                    podAffinityTerm:
                                      description: Required. A pod affinity term,
                                        associated with the corresponding weight.
                                      properties:
                                        labelSelector:
                                          description: A label query over a set of
                                            resources, in this case pods.
                                          properties:
                                            matchExpressions:
                                              description: matchExpressions is a list
                                                of label selector requirements. The
                                                requirements are ANDed.
                                              items:
                                                description: A label selector requirement
                                                  is a selector that contains values,
                                                  a key, and an operator that relates
                                                  the key and values.
                                                properties:
                                                  key:
                                                    description: key is the label
                                                      key that the selector applies
                                                      to.
                                                    type: string
                                                  operator:
                                                    description: operator represents
                                                      a key's relationship to a set
                                                      of values. Valid operators are
                                                      In, NotIn, Exists and DoesNotExist.
                                                    type: string
                                                  values:
                                                    description: values is an array
                                                      of string values. If the operator
                                                      is In or NotIn, the values array
                                                      must be non-empty. If the operator
                                                      is Exists or DoesNotExist, the
                                                      values array must be empty.
                                                      This array is replaced during
                                                      a strategic merge patch.
                                                    items:
                                                      type: string
                                                    type: array
                                                required:
                                                - key
                                                - operator
                                                type: object
                                              type: array
                                            matchLabels:
                                              additionalProperties:
                                                type: string
                                              description: matchLabels is a map of
                                                {key,value} pairs. A single {key,value}
                                                in the matchLabels map is equivalent
                                                to an element of matchExpressions,
                                                whose key field is "key", the operator
                                                is "In", and the values array contains
                                                only "value". The requirements are
                                                ANDed.
                                              type: object
                                          type: object
                                        namespaces:
                                          description: namespaces specifies which
                                            namespaces the labelSelector applies to
                                            (matches against); null or empty list
                                            means "this pod's namespace"
                                          items:
                                            type: string
                                          type: array
                                        topologyKey:
                                          description: This pod should be co-located
                                            (affinity) or not co-located (anti-affinity)
                                            with the pods matching the labelSelector
                                            in the specified namespaces, where co-located
                                            is defined as running on a node whose
                                            value of the label with key topologyKey
                                            matches that of any node on which any
                                            of the selected pods is running. Empty
                                            topologyKey is not allowed.
                                          type: string
                                      required:
                                      - topologyKey
                                      type: object
                                    weight:
                                      description: weight associated with matching
                                        the corresponding podAffinityTerm, in the
                                        range 1-100.
                                      format: int32
                                      type: integer
                                  required:
                                  - podAffinityTerm
                                  - weight
                                  type: object
                                type: array
                              requiredDuringSchedulingIgnoredDuringExecution:
                                description: If the anti-affinity requirements specified
                                  by this field are not met at scheduling time, the
                                  pod will not be scheduled onto the node. If the
                                  anti-affinity requirements specified by this field
                                  cease to be met at some point during pod execution
                                  (e.g. due to a pod label update), the system may
                                  or may not try to eventually evict the pod from
                                  its node. When there are multiple elements, the
                                  lists of nodes corresponding to each podAffinityTerm
                                  are intersected, i.e. all terms must be satisfied.
                                items:
                                  description: Defines a set of pods (namely those
                                    matching the labelSelector relative to the given
                                    namespace(s)) that this pod should be co-located
                                    (affinity) or not co-located (anti-affinity) with,
                                    where co-located is defined as running on a node
                                    whose value of the label with key <topologyKey>
                                    matches that of any node on which a pod of the
                                    set of pods is running
                                  properties:
                                    labelSelector:
                                      description: A label query over a set of resources,
                                        in this case pods.
                                      properties:
                                        matchExpressions:
                                          description: matchExpressions is a list
                                            of label selector requirements. The requirements
                                            are ANDed.
                                          items:
                                            description: A label selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: key is the label key
                                                  that the selector applies to.
                                                type: string
                                              operator:
                                                description: operator represents a
                                                  key's relationship to a set of values.
                                                  Valid operators are In, NotIn, Exists
                                                  and DoesNotExist.
                                                type: string
                                              values:
                                                description: values is an array of
                                                  string values. If the operator is
                                                  In or NotIn, the values array must
                                                  be non-empty. If the operator is
                                                  Exists or DoesNotExist, the values
                                                  array must be empty. This array
                                                  is replaced during a strategic merge
                                                  patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                        matchLabels:
                                          additionalProperties:
                                            type: string
                                          description: matchLabels is a map of {key,value}
                                            pairs. A single {key,value} in the matchLabels
                                            map is equivalent to an element of matchExpressions,
                                            whose key field is "key", the operator
                                            is "In", and the values array contains
                                            only "value". The requirements are ANDed.
                                          type: object
                                      type: object
                                    namespaces:
                                      description: namespaces specifies which namespaces
                                        the labelSelector applies to (matches against);
                                        null or empty list means "this pod's namespace"
                                      items:
                                        type: string
                                      type: array
                                    topologyKey:
                                      description: This pod should be co-located (affinity)
                                        or not co-located (anti-affinity) with the
                                        pods matching the labelSelector in the specified
                                        namespaces, where co-located is defined as
                                        running on a node whose value of the label
                                        with key topologyKey matches that of any node
                                        on which any of the selected pods is running.
                                        Empty topologyKey is not allowed.
                                      type: string
                                  required:
                                  - topologyKey
                                  type: object
                                type: array
                            type: object
                        type: object
                      nodeSelector:
                        additionalProperties:
                          type: string
                        description: 'nodeSelector is the node selector applied to
                          the relevant kind of pods It specifies a map of key-value
                          pairs: for the pod to be eligible to run on a node, the
                          node must have each of the indicated key-value pairs as
                          labels (it can have additional labels as well). See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector'
                        type: object
                      tolerations:
                        description: tolerations is a list of tolerations applied
                          to the relevant kind of pods See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
                          for more info. These are additional tolerations other than
                          default ones.
                        items:
                          description: The pod this Toleration is attached to tolerates
                            any taint that matches the triple <key,value,effect> using
                            the matching operator <operator>.
                          properties:
                            effect:
                              description: Effect indicates the taint effect to match.
                                Empty means match all taint effects. When specified,
                                allowed values are NoSchedule, PreferNoSchedule and
                                NoExecute.
                              type: string
                            key:
                              description: Key is the taint key that the toleration
                                applies to. Empty means match all taint keys. If the
                                key is empty, operator must be Exists; this combination
                                means to match all values and all keys.
                              type: string
                            operator:
                              description: Operator represents a key's relationship
                                to the value. Valid operators are Exists and Equal.
                                Defaults to Equal. Exists is equivalent to wildcard
                                for value, so that a pod can tolerate all taints of
                                a particular category.
                              type: string
                            tolerationSeconds:
                              description: TolerationSeconds represents the period
                                of time the toleration (which must be of effect NoExecute,
                                otherwise this field is ignored) tolerates the taint.
                                By default, it is not set, which means tolerate the
                                taint forever (do not evict). Zero and negative values
                                will be treated as 0 (evict immediately) by the system.
                              format: int64
                              type: integer
                            value:
                              description: Value is the taint value the toleration
                                matches to. If the operator is Exists, the value should
                                be empty, otherwise just a regular string.
                              type: string
                          type: object
                        type: array
                    type: object
                type: object
              localStorageClassName:
                description: LocalStorageClassName the name of the local storage class.
                type: string
              version:
                description: operator version
                type: string
              workloads:
                description: workloads HyperConvergedConfig influences the pod configuration
                  (currently only placement) of components which need to be running
                  on a node where virtualization workloads should be able to run.
                  Changes to Workloads HyperConvergedConfig can be applied only without
                  existing workload.
                properties:
                  nodePlacement:
                    description: NodePlacement describes node scheduling configuration.
                    properties:
                      affinity:
                        description: affinity enables pod affinity/anti-affinity placement
                          expanding the types of constraints that can be expressed
                          with nodeSelector. affinity is going to be applied to the
                          relevant kind of pods in parallel with nodeSelector See
                          https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#affinity-and-anti-affinity
                        properties:
                          nodeAffinity:
                            description: Describes node affinity scheduling rules
                              for the pod.
                            properties:
                              preferredDuringSchedulingIgnoredDuringExecution:
                                description: The scheduler will prefer to schedule
                                  pods to nodes that satisfy the affinity expressions
                                  specified by this field, but it may choose a node
                                  that violates one or more of the expressions. The
                                  node that is most preferred is the one with the
                                  greatest sum of weights, i.e. for each node that
                                  meets all of the scheduling requirements (resource
                                  request, requiredDuringScheduling affinity expressions,
                                  etc.), compute a sum by iterating through the elements
                                  of this field and adding "weight" to the sum if
                                  the node matches the corresponding matchExpressions;
                                  the node(s) with the highest sum are the most preferred.
                                items:
                                  description: An empty preferred scheduling term
                                    matches all objects with implicit weight 0 (i.e.
                                    it's a no-op). A null preferred scheduling term
                                    matches no objects (i.e. is also a no-op).
                                  properties:
                                    preference:
                                      description: A node selector term, associated
                                        with the corresponding weight.
                                      properties:
                                        matchExpressions:
                                          description: A list of node selector requirements
                                            by node's labels.
                                          items:
                                            description: A node selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: The label key that the
                                                  selector applies to.
                                                type: string
                                              operator:
                                                description: Represents a key's relationship
                                                  to a set of values. Valid operators
                                                  are In, NotIn, Exists, DoesNotExist.
                                                  Gt, and Lt.
                                                type: string
                                              values:
                                                description: An array of string values.
                                                  If the operator is In or NotIn,
                                                  the values array must be non-empty.
                                                  If the operator is Exists or DoesNotExist,
                                                  the values array must be empty.
                                                  If the operator is Gt or Lt, the
                                                  values array must have a single
                                                  element, which will be interpreted
                                                  as an integer. This array is replaced
                                                  during a strategic merge patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                        matchFields:
                                          description: A list of node selector requirements
                                            by node's fields.
                                          items:
                                            description: A node selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: The label key that the
                                                  selector applies to.
                                                type: string
                                              operator:
                                                description: Represents a key's relationship
                                                  to a set of values. Valid operators
                                                  are In, NotIn, Exists, DoesNotExist.
                                                  Gt, and Lt.
                                                type: string
                                              values:
                                                description: An array of string values.
                                                  If the operator is In or NotIn,
                                                  the values array must be non-empty.
                                                  If the operator is Exists or DoesNotExist,
                                                  the values array must be empty.
                                                  If the operator is Gt or Lt, the
                                                  values array must have a single
                                                  element, which will be interpreted
                                                  as an integer. This array is replaced
                                                  during a strategic merge patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                      type: object
                                    weight:
                                      description: Weight associated with matching
                                        the corresponding nodeSelectorTerm, in the
                                        range 1-100.
                                      format: int32
                                      type: integer
                                  required:
                                  - preference
                                  - weight
                                  type: object
                                type: array
                              requiredDuringSchedulingIgnoredDuringExecution:
                                description: If the affinity requirements specified
                                  by this field are not met at scheduling time, the
                                  pod will not be scheduled onto the node. If the
                                  affinity requirements specified by this field cease
                                  to be met at some point during pod execution (e.g.
                                  due to an update), the system may or may not try
                                  to eventually evict the pod from its node.
                                properties:
                                  nodeSelectorTerms:
                                    description: Required. A list of node selector
                                      terms. The terms are ORed.
                                    items:
                                      description: A null or empty node selector term
                                        matches no objects. The requirements of them
                                        are ANDed. The TopologySelectorTerm type implements
                                        a subset of the NodeSelectorTerm.
                                      properties:
                                        matchExpressions:
                                          description: A list of node selector requirements
                                            by node's labels.
                                          items:
                                            description: A node selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: The label key that the
                                                  selector applies to.
                                                type: string
                                              operator:
                                                description: Represents a key's relationship
                                                  to a set of values. Valid operators
                                                  are In, NotIn, Exists, DoesNotExist.
                                                  Gt, and Lt.
                                                type: string
                                              values:
                                                description: An array of string values.
                                                  If the operator is In or NotIn,
                                                  the values array must be non-empty.
                                                  If the operator is Exists or DoesNotExist,
                                                  the values array must be empty.
                                                  If the operator is Gt or Lt, the
                                                  values array must have a single
                                                  element, which will be interpreted
                                                  as an integer. This array is replaced
                                                  during a strategic merge patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                        matchFields:
                                          description: A list of node selector requirements
                                            by node's fields.
                                          items:
                                            description: A node selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: The label key that the
                                                  selector applies to.
                                                type: string
                                              operator:
                                                description: Represents a key's relationship
                                                  to a set of values. Valid operators
                                                  are In, NotIn, Exists, DoesNotExist.
                                                  Gt, and Lt.
                                                type: string
                                              values:
                                                description: An array of string values.
                                                  If the operator is In or NotIn,
                                                  the values array must be non-empty.
                                                  If the operator is Exists or DoesNotExist,
                                                  the values array must be empty.
                                                  If the operator is Gt or Lt, the
                                                  values array must have a single
                                                  element, which will be interpreted
                                                  as an integer. This array is replaced
                                                  during a strategic merge patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                      type: object
                                    type: array
                                required:
                                - nodeSelectorTerms
                                type: object
                            type: object
                          podAffinity:
                            description: Describes pod affinity scheduling rules (e.g.
                              co-locate this pod in the same node, zone, etc. as some
                              other pod(s)).
                            properties:
                              preferredDuringSchedulingIgnoredDuringExecution:
                                description: The scheduler will prefer to schedule
                                  pods to nodes that satisfy the affinity expressions
                                  specified by this field, but it may choose a node
                                  that violates one or more of the expressions. The
                                  node that is most preferred is the one with the
                                  greatest sum of weights, i.e. for each node that
                                  meets all of the scheduling requirements (resource
                                  request, requiredDuringScheduling affinity expressions,
                                  etc.), compute a sum by iterating through the elements
                                  of this field and adding "weight" to the sum if
                                  the node has pods which matches the corresponding
                                  podAffinityTerm; the node(s) with the highest sum
                                  are the most preferred.
                                items:
                                  description: The weights of all of the matched WeightedPodAffinityTerm
                                    fields are added per-node to find the most preferred
                                    node(s)
                                  properties:
                                    podAffinityTerm:
                                      description: Required. A pod affinity term,
                                        associated with the corresponding weight.
                                      properties:
                                        labelSelector:
                                          description: A label query over a set of
                                            resources, in this case pods.
                                          properties:
                                            matchExpressions:
                                              description: matchExpressions is a list
                                                of label selector requirements. The
                                                requirements are ANDed.
                                              items:
                                                description: A label selector requirement
                                                  is a selector that contains values,
                                                  a key, and an operator that relates
                                                  the key and values.
                                                properties:
                                                  key:
                                                    description: key is the label
                                                      key that the selector applies
                                                      to.
                                                    type: string
                                                  operator:
                                                    description: operator represents
                                                      a key's relationship to a set
                                                      of values. Valid operators are
                                                      In, NotIn, Exists and DoesNotExist.
                                                    type: string
                                                  values:
                                                    description: values is an array
                                                      of string values. If the operator
                                                      is In or NotIn, the values array
                                                      must be non-empty. If the operator
                                                      is Exists or DoesNotExist, the
                                                      values array must be empty.
                                                      This array is replaced during
                                                      a strategic merge patch.
                                                    items:
                                                      type: string
                                                    type: array
                                                required:
                                                - key
                                                - operator
                                                type: object
                                              type: array
                                            matchLabels:
                                              additionalProperties:
                                                type: string
                                              description: matchLabels is a map of
                                                {key,value} pairs. A single {key,value}
                                                in the matchLabels map is equivalent
                                                to an element of matchExpressions,
                                                whose key field is "key", the operator
                                                is "In", and the values array contains
                                                only "value". The requirements are
                                                ANDed.
                                              type: object
                                          type: object
                                        namespaces:
                                          description: namespaces specifies which
                                            namespaces the labelSelector applies to
                                            (matches against); null or empty list
                                            means "this pod's namespace"
                                          items:
                                            type: string
                                          type: array
                                        topologyKey:
                                          description: This pod should be co-located
                                            (affinity) or not co-located (anti-affinity)
                                            with the pods matching the labelSelector
                                            in the specified namespaces, where co-located
                                            is defined as running on a node whose
                                            value of the label with key topologyKey
                                            matches that of any node on which any
                                            of the selected pods is running. Empty
                                            topologyKey is not allowed.
                                          type: string
                                      required:
                                      - topologyKey
                                      type: object
                                    weight:
                                      description: weight associated with matching
                                        the corresponding podAffinityTerm, in the
                                        range 1-100.
                                      format: int32
                                      type: integer
                                  required:
                                  - podAffinityTerm
                                  - weight
                                  type: object
                                type: array
                              requiredDuringSchedulingIgnoredDuringExecution:
                                description: If the affinity requirements specified
                                  by this field are not met at scheduling time, the
                                  pod will not be scheduled onto the node. If the
                                  affinity requirements specified by this field cease
                                  to be met at some point during pod execution (e.g.
                                  due to a pod label update), the system may or may
                                  not try to eventually evict the pod from its node.
                                  When there are multiple elements, the lists of nodes
                                  corresponding to each podAffinityTerm are intersected,
                                  i.e. all terms must be satisfied.
                                items:
                                  description: Defines a set of pods (namely those
                                    matching the labelSelector relative to the given
                                    namespace(s)) that this pod should be co-located
                                    (affinity) or not co-located (anti-affinity) with,
                                    where co-located is defined as running on a node
                                    whose value of the label with key <topologyKey>
                                    matches that of any node on which a pod of the
                                    set of pods is running
                                  properties:
                                    labelSelector:
                                      description: A label query over a set of resources,
                                        in this case pods.
                                      properties:
                                        matchExpressions:
                                          description: matchExpressions is a list
                                            of label selector requirements. The requirements
                                            are ANDed.
                                          items:
                                            description: A label selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: key is the label key
                                                  that the selector applies to.
                                                type: string
                                              operator:
                                                description: operator represents a
                                                  key's relationship to a set of values.
                                                  Valid operators are In, NotIn, Exists
                                                  and DoesNotExist.
                                                type: string
                                              values:
                                                description: values is an array of
                                                  string values. If the operator is
                                                  In or NotIn, the values array must
                                                  be non-empty. If the operator is
                                                  Exists or DoesNotExist, the values
                                                  array must be empty. This array
                                                  is replaced during a strategic merge
                                                  patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                        matchLabels:
                                          additionalProperties:
                                            type: string
                                          description: matchLabels is a map of {key,value}
                                            pairs. A single {key,value} in the matchLabels
                                            map is equivalent to an element of matchExpressions,
                                            whose key field is "key", the operator
                                            is "In", and the values array contains
                                            only "value". The requirements are ANDed.
                                          type: object
                                      type: object
                                    namespaces:
                                      description: namespaces specifies which namespaces
                                        the labelSelector applies to (matches against);
                                        null or empty list means "this pod's namespace"
                                      items:
                                        type: string
                                      type: array
                                    topologyKey:
                                      description: This pod should be co-located (affinity)
                                        or not co-located (anti-affinity) with the
                                        pods matching the labelSelector in the specified
                                        namespaces, where co-located is defined as
                                        running on a node whose value of the label
                                        with key topologyKey matches that of any node
                                        on which any of the selected pods is running.
                                        Empty topologyKey is not allowed.
                                      type: string
                                  required:
                                  - topologyKey
                                  type: object
                                type: array
                            type: object
                          podAntiAffinity:
                            description: Describes pod anti-affinity scheduling rules
                              (e.g. avoid putting this pod in the same node, zone,
                              etc. as some other pod(s)).
                            properties:
                              preferredDuringSchedulingIgnoredDuringExecution:
                                description: The scheduler will prefer to schedule
                                  pods to nodes that satisfy the anti-affinity expressions
                                  specified by this field, but it may choose a node
                                  that violates one or more of the expressions. The
                                  node that is most preferred is the one with the
                                  greatest sum of weights, i.e. for each node that
                                  meets all of the scheduling requirements (resource
                                  request, requiredDuringScheduling anti-affinity
                                  expressions, etc.), compute a sum by iterating through
                                  the elements of this field and adding "weight" to
                                  the sum if the node has pods which matches the corresponding
                                  podAffinityTerm; the node(s) with the highest sum
                                  are the most preferred.
                                items:
                                  description: The weights of all of the matched WeightedPodAffinityTerm
                                    fields are added per-node to find the most preferred
                                    node(s)
                                  properties:
                                    podAffinityTerm:
                                      description: Required. A pod affinity term,
                                        associated with the corresponding weight.
                                      properties:
                                        labelSelector:
                                          description: A label query over a set of
                                            resources, in this case pods.
                                          properties:
                                            matchExpressions:
                                              description: matchExpressions is a list
                                                of label selector requirements. The
                                                requirements are ANDed.
                                              items:
                                                description: A label selector requirement
                                                  is a selector that contains values,
                                                  a key, and an operator that relates
                                                  the key and values.
                                                properties:
                                                  key:
                                                    description: key is the label
                                                      key that the selector applies
                                                      to.
                                                    type: string
                                                  operator:
                                                    description: operator represents
                                                      a key's relationship to a set
                                                      of values. Valid operators are
                                                      In, NotIn, Exists and DoesNotExist.
                                                    type: string
                                                  values:
                                                    description: values is an array
                                                      of string values. If the operator
                                                      is In or NotIn, the values array
                                                      must be non-empty. If the operator
                                                      is Exists or DoesNotExist, the
                                                      values array must be empty.
                                                      This array is replaced during
                                                      a strategic merge patch.
                                                    items:
                                                      type: string
                                                    type: array
                                                required:
                                                - key
                                                - operator
                                                type: object
                                              type: array
                                            matchLabels:
                                              additionalProperties:
                                                type: string
                                              description: matchLabels is a map of
                                                {key,value} pairs. A single {key,value}
                                                in the matchLabels map is equivalent
                                                to an element of matchExpressions,
                                                whose key field is "key", the operator
                                                is "In", and the values array contains
                                                only "value". The requirements are
                                                ANDed.
                                              type: object
                                          type: object
                                        namespaces:
                                          description: namespaces specifies which
                                            namespaces the labelSelector applies to
                                            (matches against); null or empty list
                                            means "this pod's namespace"
                                          items:
                                            type: string
                                          type: array
                                        topologyKey:
                                          description: This pod should be co-located
                                            (affinity) or not co-located (anti-affinity)
                                            with the pods matching the labelSelector
                                            in the specified namespaces, where co-located
                                            is defined as running on a node whose
                                            value of the label with key topologyKey
                                            matches that of any node on which any
                                            of the selected pods is running. Empty
                                            topologyKey is not allowed.
                                          type: string
                                      required:
                                      - topologyKey
                                      type: object
                                    weight:
                                      description: weight associated with matching
                                        the corresponding podAffinityTerm, in the
                                        range 1-100.
                                      format: int32
                                      type: integer
                                  required:
                                  - podAffinityTerm
                                  - weight
                                  type: object
                                type: array
                              requiredDuringSchedulingIgnoredDuringExecution:
                                description: If the anti-affinity requirements specified
                                  by this field are not met at scheduling time, the
                                  pod will not be scheduled onto the node. If the
                                  anti-affinity requirements specified by this field
                                  cease to be met at some point during pod execution
                                  (e.g. due to a pod label update), the system may
                                  or may not try to eventually evict the pod from
                                  its node. When there are multiple elements, the
                                  lists of nodes corresponding to each podAffinityTerm
                                  are intersected, i.e. all terms must be satisfied.
                                items:
                                  description: Defines a set of pods (namely those
                                    matching the labelSelector relative to the given
                                    namespace(s)) that this pod should be co-located
                                    (affinity) or not co-located (anti-affinity) with,
                                    where co-located is defined as running on a node
                                    whose value of the label with key <topologyKey>
                                    matches that of any node on which a pod of the
                                    set of pods is running
                                  properties:
                                    labelSelector:
                                      description: A label query over a set of resources,
                                        in this case pods.
                                      properties:
                                        matchExpressions:
                                          description: matchExpressions is a list
                                            of label selector requirements. The requirements
                                            are ANDed.
                                          items:
                                            description: A label selector requirement
                                              is a selector that contains values,
                                              a key, and an operator that relates
                                              the key and values.
                                            properties:
                                              key:
                                                description: key is the label key
                                                  that the selector applies to.
                                                type: string
                                              operator:
                                                description: operator represents a
                                                  key's relationship to a set of values.
                                                  Valid operators are In, NotIn, Exists
                                                  and DoesNotExist.
                                                type: string
                                              values:
                                                description: values is an array of
                                                  string values. If the operator is
                                                  In or NotIn, the values array must
                                                  be non-empty. If the operator is
                                                  Exists or DoesNotExist, the values
                                                  array must be empty. This array
                                                  is replaced during a strategic merge
                                                  patch.
                                                items:
                                                  type: string
                                                type: array
                                            required:
                                            - key
                                            - operator
                                            type: object
                                          type: array
                                        matchLabels:
                                          additionalProperties:
                                            type: string
                                          description: matchLabels is a map of {key,value}
                                            pairs. A single {key,value} in the matchLabels
                                            map is equivalent to an element of matchExpressions,
                                            whose key field is "key", the operator
                                            is "In", and the values array contains
                                            only "value". The requirements are ANDed.
                                          type: object
                                      type: object
                                    namespaces:
                                      description: namespaces specifies which namespaces
                                        the labelSelector applies to (matches against);
                                        null or empty list means "this pod's namespace"
                                      items:
                                        type: string
                                      type: array
                                    topologyKey:
                                      description: This pod should be co-located (affinity)
                                        or not co-located (anti-affinity) with the
                                        pods matching the labelSelector in the specified
                                        namespaces, where co-located is defined as
                                        running on a node whose value of the label
                                        with key topologyKey matches that of any node
                                        on which any of the selected pods is running.
                                        Empty topologyKey is not allowed.
                                      type: string
                                  required:
                                  - topologyKey
                                  type: object
                                type: array
                            type: object
                        type: object
                      nodeSelector:
                        additionalProperties:
                          type: string
                        description: 'nodeSelector is the node selector applied to
                          the relevant kind of pods It specifies a map of key-value
                          pairs: for the pod to be eligible to run on a node, the
                          node must have each of the indicated key-value pairs as
                          labels (it can have additional labels as well). See https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector'
                        type: object
                      tolerations:
                        description: tolerations is a list of tolerations applied
                          to the relevant kind of pods See https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/
                          for more info. These are additional tolerations other than
                          default ones.
                        items:
                          description: The pod this Toleration is attached to tolerates
                            any taint that matches the triple <key,value,effect> using
                            the matching operator <operator>.
                          properties:
                            effect:
                              description: Effect indicates the taint effect to match.
                                Empty means match all taint effects. When specified,
                                allowed values are NoSchedule, PreferNoSchedule and
                                NoExecute.
                              type: string
                            key:
                              description: Key is the taint key that the toleration
                                applies to. Empty means match all taint keys. If the
                                key is empty, operator must be Exists; this combination
                                means to match all values and all keys.
                              type: string
                            operator:
                              description: Operator represents a key's relationship
                                to the value. Valid operators are Exists and Equal.
                                Defaults to Equal. Exists is equivalent to wildcard
                                for value, so that a pod can tolerate all taints of
                                a particular category.
                              type: string
                            tolerationSeconds:
                              description: TolerationSeconds represents the period
                                of time the toleration (which must be of effect NoExecute,
                                otherwise this field is ignored) tolerates the taint.
                                By default, it is not set, which means tolerate the
                                taint forever (do not evict). Zero and negative values
                                will be treated as 0 (evict immediately) by the system.
                              format: int64
                              type: integer
                            value:
                              description: Value is the taint value the toleration
                                matches to. If the operator is Exists, the value should
                                be empty, otherwise just a regular string.
                              type: string
                          type: object
                        type: array
                    type: object
                type: object
            type: object
          status:
            description: HyperConvergedStatus defines the observed state of HyperConverged
            properties:
              conditions:
                description: Conditions describes the state of the HyperConverged
                  resource.
                items:
                  description: Condition represents the state of the operator's reconciliation
                    functionality.
                  properties:
                    lastHeartbeatTime:
                      format: date-time
                      type: string
                    lastTransitionTime:
                      format: date-time
                      type: string
                    message:
                      type: string
                    reason:
                      type: string
                    status:
                      type: string
                    type:
                      description: ConditionType is the state of the operator's reconciliation
                        functionality.
                      type: string
                  required:
                  - status
                  - type
                  type: object
                type: array
              relatedObjects:
                description: RelatedObjects is a list of objects created and maintained
                  by this operator. Object references will be added to this list after
                  they have been created AND found in the cluster.
                items:
                  description: 'ObjectReference contains enough information to let
                    you inspect or modify the referred object. --- New uses of this
                    type are discouraged because of difficulty describing its usage
                    when embedded in APIs.  1. Ignored fields.  It includes many fields
                    which are not generally honored.  For instance, ResourceVersion
                    and FieldPath are both very rarely valid in actual usage.  2.
                    Invalid usage help.  It is impossible to add specific help for
                    individual usage.  In most embedded usages, there are particular     restrictions
                    like, "must refer only to types A and B" or "UID not honored"
                    or "name must be restricted".     Those cannot be well described
                    when embedded.  3. Inconsistent validation.  Because the usages
                    are different, the validation rules are different by usage, which
                    makes it hard for users to predict what will happen.  4. The fields
                    are both imprecise and overly precise.  Kind is not a precise
                    mapping to a URL. This can produce ambiguity     during interpretation
                    and require a REST mapping.  In most cases, the dependency is
                    on the group,resource tuple     and the version of the actual
                    struct is irrelevant.  5. We cannot easily change it.  Because
                    this type is embedded in many locations, updates to this type     will
                    affect numerous schemas.  Don''t make new APIs embed an underspecified
                    API type they do not control. Instead of using this type, create
                    a locally provided and used type that is well-focused on your
                    reference. For example, ServiceReferences for admission registration:
                    https://github.com/kubernetes/api/blob/release-1.17/admissionregistration/v1/types.go#L533
                    .'
                  properties:
                    apiVersion:
                      description: API version of the referent.
                      type: string
                    fieldPath:
                      description: 'If referring to a piece of an object instead of
                        an entire object, this string should contain a valid JSON/Go
                        field access statement, such as desiredState.manifest.containers[2].
                        For example, if the object reference is to a container within
                        a pod, this would take on a value like: "spec.containers{name}"
                        (where "name" refers to the name of the container that triggered
                        the event) or if no container name is specified "spec.containers[2]"
                        (container with index 2 in this pod). This syntax is chosen
                        only to have some well-defined way of referencing a part of
                        an object. TODO: this design is not final and this field is
                        subject to change in the future.'
                      type: string
                    kind:
                      description: 'Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
                      type: string
                    name:
                      description: 'Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names'
                      type: string
                    namespace:
                      description: 'Namespace of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/'
                      type: string
                    resourceVersion:
                      description: 'Specific resourceVersion to which this reference
                        is made, if any. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency'
                      type: string
                    uid:
                      description: 'UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids'
                      type: string
                  type: object
                type: array
              versions:
                description: 'Versions is a list of HCO component versions, as name/version
                  pairs. The version with a name of "operator" is the HCO version
                  itself, as described here: https://github.com/openshift/cluster-version-operator/blob/master/docs/dev/clusteroperator.md#version'
                items:
                  properties:
                    name:
                      type: string
                    version:
                      type: string
                  type: object
                type: array
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}`
