package docs

// get pods memory
var DocsGetPodsMemory = map[string]string{
	"hostName": "The name of the Kubernetes node where the pod is running. Nodes are the physical or virtual machines that form the cluster and run your workloads.",

	"podName": "The name of the pod. A pod is the smallest deployable unit in Kubernetes and can contain one or more containers that share resources and a network namespace.",

	"usageBytes": "The total memory used by the pod or container, including both active and cached memory. This metric provides a comprehensive view of the memory footprint, but it may include memory that can be reclaimed by the OS (like caches and buffers).",

	"workingSetBytes": `The amount of memory actively used by the pod that cannot be reclaimed by the OS. This excludes memory that is cached and can be evicted under memory pressure. 
The pod's workingSetBytes is the total of all containers' working set memory, representing the minimum memory required by the pod to function.`,

	"container": "The name of the container within the pod. Containers are individual units of application execution within a pod, each with its own isolated filesystem and process space.",

	"containerWorkingSetBytes": `The memory actively in use by this container. It excludes caches and buffers that can be reclaimed by the OS. 
This value may differ from the pod-level workingSetBytes because the pod's memory usage is aggregated across all containers, and shared memory between containers may reduce the overall pod memory footprint.`,

	"Hint": `In some cases, the workingSetBytes of individual containers may add up to more than the pod's workingSetBytes. This can happen because containers might share memory (such as shared libraries or system files) that are counted towards each container's working set, but not double-counted at the pod level. 
This memory optimization reduces the overall memory footprint for the pod while showing higher individual working set memory for the containers.`,

	"Hint 2": `At the OS and hardware level, the actual memory consumption is based on the pod's workingSetBytes, not the sum of the containers' workingSetBytes. 
This is because shared memory (such as libraries or system files) is only counted once at the pod level, even if it is counted separately for each container. 
Thus, the operating system and hardware allocate memory according to the pod's workingSetBytes, reflecting the true memory usage on the system.`,
}
