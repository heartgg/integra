import { Button, StyleSheet, Text, View, ActivityIndicator } from "react-native";
import DropDownPicker from "react-native-dropdown-picker";
import React, { useState, useEffect } from "react";
import { BarCodeScanner } from "expo-barcode-scanner";
import SuccessAnimation from "./SuccessAnimation.js";
import * as Location from "expo-location";

export default function App() {
  const [hasCameraPermission, setHasCameraPermission] = useState(null);
  const [hasLocationPermission, setHasLocationPermission] = useState(null);
  const [scanned, setScanned] = useState(false);
  const [scanTimeout, setScanTimeout] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);
  const [patientID, setPatientID] = useState(null);
  const [open, setOpen] = useState(false);
  const [modality, setModality] = useState(null);
  const [modalityList, setModalityList] = useState([
    { label: "IE Fluoro", value: "Fluoro" },
    { label: "XRAY", value: "XRAY" },
    { label: "CT", value: "CT" },
    { label: "IR", value: "IR" },
    { label: "MRI", value: "MRI" },
    { label: "US", value: "US" },
    { label: "Dexa", value: "Dexa" },
    { label: "Nuc Med", value: "NucMed" },
  ]);

  useEffect(() => {
    (async () => {
      let { status } = await BarCodeScanner.requestPermissionsAsync();
      setHasCameraPermission(status === "granted");
      status = await Location.requestForegroundPermissionsAsync();
      setHasLocationPermission(status === "granted");
    })();
  }, []);

  const handleBarCodeScanned = async ({ type, data }) => {
    if (scanned) return;
    if (!modality) {
      alert("Please select a modality before scanning.");
      setScanTimeout(
        setTimeout(() => {
          setScanTimeout(null);
        }, 3000)
      );
      return;
    }

    setIsLoading(true);
    setScanned(true);
    try {
      const location = await Location.getCurrentPositionAsync({
        accuracy: 3,
      });
      console.log(location);
      setPatientID(data);
      console.log(`https://integri-scan.herokuapp.com/scan-exams
          ?modality=${modality}
          &patientID=${data}
          &latitude=${location.coords.latitude}
          &longitude=${location.coords.longitude}`);
      //send message to server with scanned patient ID
      await fetch(
        `https://integri-scan.herokuapp.com/scan-exams?modality=${modality}&patientID=${data}&latitude=${location.coords.latitude}&longitude=${location.coords.longitude}`
      );
      setIsLoading(false);
      setShowSuccess(true);
    } catch (err) {
      alert(
        "An error occurred when trying to send patient data. Please try again."
      );
      setIsLoading(false);
    } 
  };

  return (
    <View style={styles.container}>
      <View style={styles.infoText}>
        <Text>Please scan patient barcode</Text>
      </View>
      {hasCameraPermission == null ? (
        <Text>Requesting for camera permission</Text>
      ) : hasCameraPermission == false ? (
        <Text>No access to camera. Please go to app permission settings and allow camera permissions.</Text>
      ) : (
        <BarCodeScanner
          onBarCodeScanned={scanTimeout != null ? undefined : handleBarCodeScanned}
          style={StyleSheet.absoluteFillObject}
        />
      )}
      {(scanned) && (
        <View style={styles.scanAgainButton}>
          <Button
            title={"Tap to Scan Again"}
            onPress={() => setScanned(false)}
          />
        </View>
      )}
      {isLoading && <ActivityIndicator color={"00f"} size="large"/>}
      {showSuccess && (
        <SuccessAnimation
          onAnimationEnd={() => {
            alert(`Successfully scanned patient ID: ${patientID}`);
            setShowSuccess(false);
          }}
        />
      )}
      <View style={styles.selectModality}>
        <DropDownPicker
          open={open}
          value={modality}
          items={modalityList}
          setOpen={setOpen}
          setValue={setModality}
          setItems={setModalityList}
        />
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    alignItems: "center",
    justifyContent: "center",
  },
  infoText: {
    flex: 1,
    alignItems: "center",
    justifyContent: "flex-start",
    marginTop: 30,
  },
  selectModality: {
    flex: 1,
    alignItems: "center",
    justifyContent: "flex-end",
    marginBottom: 3,
  },
  scanAgainButton: {
    flex: 1,
    alignItems: "center",
    justifyContent: "flex-start",
  },
  barCodeView: {
    width: "100%",
    height: "50%",
    marginBottom: 30,
  },
});
