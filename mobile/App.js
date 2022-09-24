import { Button, StyleSheet, Text, View } from "react-native";
import DropDownPicker from "react-native-dropdown-picker";
import React, { useState, useEffect } from "react";
import { BarCodeScanner } from "expo-barcode-scanner";
import SuccessAnimation from "./SuccessAnimation.js";

export default function App() {
  // camera permissions
  const [hasPermission, setHasPermission] = useState(null);
  const [scanned, setScanned] = useState(false);
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
      const { status } = await BarCodeScanner.requestPermissionsAsync();
      setHasPermission(status === "granted");
    })();
  }, []);

  const handleBarCodeScanned = async ({ type, data }) => {
    setScanned(true);
    setPatientID(data);
    // send message to server with scanned patient ID
    try {
      await fetch(
        `https://integri-scan.herokuapp.com/scan-exams?modality=${modality}&patientID=${data}`
      );
      setShowSuccess(true);
    } catch (err) {
      alert(
        "An error occurred when trying to send patient data. Please try again."
      );
    }
  };

  if (hasPermission === null) {
    return <Text>Requesting for camera permission</Text>;
  }
  if (hasPermission === false) {
    return <Text>No access to camera</Text>;
  }

  return (
    <View style={styles.container}>
      <View style={styles.infoText}>
        <Text>Please scan patient barcode</Text>
      </View>
      <BarCodeScanner
        onBarCodeScanned={
          scanned || !modality ? undefined : handleBarCodeScanned
        }
        style={StyleSheet.absoluteFillObject}
      />
      {showSuccess && (
        <SuccessAnimation
          onAnimationEnd={() => {
            alert(`Successfully scanned patient ID: ${patientID}`);
            setShowSuccess(false);
          }}
        />
      )}
      <View style={styles.selectModality}>
        {scanned && (
          <Button
            title={"Tap to Scan Again"}
            onPress={() => setScanned(false)}
          />
        )}
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
  },
  barCodeView: {
    width: "100%",
    height: "50%",
    marginBottom: 30,
  },
});
