/*
 * Copyright (c) 2020. Ontario Institute for Cancer Research
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"

	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"net/http"
	"strconv"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var (
	podResource    = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}

	// Defaults
	overrideVolumePathCollision = true
	scratchDirName = "/icgc-argo-scratch"
	scratchVolumeName = "icgc-argo-scratch"
	targetNamespace = ""
	debug = false
)

type EmptyDirData struct {
	Name string `json:"name"`
	EmptyDir interface{} `json:"emptyDir"`
}

func buildEmptyDirDatas() []EmptyDirData{
	var emptyDirData = EmptyDirData{ Name: scratchVolumeName, EmptyDir: struct {}{}}
	var emptyDirDatas []EmptyDirData
	return append(emptyDirDatas, emptyDirData)
}

func buildEmptyDirVolumeMounts() []corev1.VolumeMount{
	var volumeMount = corev1.VolumeMount{Name: scratchVolumeName, MountPath: scratchDirName}
	var volumeMounts []corev1.VolumeMount
	return append(volumeMounts, volumeMount)
}

// Adds the correct Json Patch to the patches variable for the volumes section
func appendEmptyDirPatch(patches []patchOperation) []patchOperation {
	var emptyDirDatas = buildEmptyDirDatas()
	patches = append(patches, patchOperation{
		Op:    "add",
		Path:  "/spec/volumes",
		Value: emptyDirDatas,
	})
	return patches
}

// Adds the correct Json Patch to the patches variable for the container volume mounts section
func appendVolumeMountPatch(patches []patchOperation, containerPos int, volumeMountPos int,
	containerVolumeMount *corev1.VolumeMount) []patchOperation{
	var emptyDirVolumeMounts = buildEmptyDirVolumeMounts()

	if containerVolumeMount == nil{
		patches = append(patches, patchOperation{
			Op:    "add",
			Path:  "/spec/containers/"+strconv.Itoa(containerPos)+"/volumeMounts",
			Value: emptyDirVolumeMounts,
		})

	} else {
		if overrideVolumePathCollision{
			log.Println("Container volume mount ",scratchVolumeName," already exists but overriding ")
			patches = append(patches, patchOperation{
				Op:    "replace",
				Path:  "/spec/containers/"+strconv.Itoa(containerPos)+"/volumeMounts/"+strconv.Itoa(volumeMountPos),
				Value: emptyDirVolumeMounts,
			})
		} else {
			log.Println("Container volume mount ",scratchVolumeName," already exists, and NOT overriding ")
		}
	}
	return patches
}

func dumpPodSpecs(pod *corev1.Pod){
	out, err := json.Marshal(pod)
	if  err == nil {
		log.Println("Dump Pod spec before mutation: ", string(out))
		log.Println("********************************************************")
	} else {
		log.Println("ERROR DUMPING POD SPEC: ", err)
	}

}

func dumpPatches(patches []patchOperation) {
	out, err := json.Marshal(patches)
	if  err == nil {
		log.Println("Dump Patches: ", string(out))
		log.Println("********************************************************")
	} else {
		log.Println("ERROR DUMPING PATCHES: ", err)
	}
}

// Core mutation logic
func applySecurityDefaults(req *v1beta1.AdmissionRequest) ([]patchOperation, error) {
	var patches []patchOperation

	var pod, err = extractPodSpec(req)
	if err != nil {
		return patches, err
	}

	if debug {
		dumpPodSpecs(&pod)
	}

	if !isPodInNamespace(&pod, targetNamespace){
		log.Printf("Pod request with name '%s' does not belong to targetNamespace '%s', skipping mutation.\n", pod.Name, targetNamespace)
		return patches, nil
	} else {
		log.Printf("Pod request with name '%s' detected for targetNamespace '%s'. Continuing with mutation.\n", pod.Name, targetNamespace)
	}

	if hasVolume(&pod, scratchVolumeName){
		log.Println("Already contains the scratch volume name: ", scratchVolumeName)
		return patches, nil
	}

	patches = appendEmptyDirPatch(patches)

	if pod.Spec.Containers != nil {
		for containerPos, container := range pod.Spec.Containers {
			var containerVolumeMount, volumeMountPos = findVolumeMount(&container)
			patches =  appendVolumeMountPatch(patches, containerPos, volumeMountPos, containerVolumeMount)
		}
	}

	if debug {
		//dumpPodSpecs(&pod)
		dumpPatches(patches)
	}
	return patches, nil
}

func main() {
	// Bind the config to the variables
	var cfg = parseConfig()
	overrideVolumePathCollision = cfg.App.OverrideVolumeCollisions
	scratchDirName = cfg.App.EmptyDir.MountPath
	scratchVolumeName = cfg.App.EmptyDir.VolumeName
	targetNamespace = cfg.App.TargetNamespace
	debug = cfg.App.Debug

	// Start server
	mux := http.NewServeMux()
	mux.Handle("/mutate", admitFuncHandler(applySecurityDefaults))
	mux.Handle("/health", healthFuncHandler())
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    ":"+cfg.Server.Port,
		Handler: mux,
	}

	log.Println("Starting server on port ",cfg.Server.Port," with TLS ENABLED=",cfg.Server.SSL.Enable)
	if cfg.Server.SSL.Enable {
		log.Fatal(server.ListenAndServeTLS(cfg.Server.SSL.CertPath, cfg.Server.SSL.KeyPath))
	} else {
		log.Fatal(server.ListenAndServe())
	}
	log.Println("Stopped server")
}

