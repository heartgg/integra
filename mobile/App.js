import { Button, StyleSheet, Text, View } from "react-native";
import DropDownPicker from "react-native-dropdown-picker";
import React, { useState, useEffect } from "react";
import { BarCodeScanner } from "expo-barcode-scanner";

export default function App() {
  // camera permissions
  const [hasPermission, setHasPermission] = useState(null);
  const [scanned, setScanned] = useState(false);
  const [open, setOpen] = useState(false);
  const [modality, setModality] = useState(null);
  const [modalityList, setModalityList] = useState([
    { label: "IE Fluoro", value: "IE Fluoro" },
    { label: "XRAY", value: "XRAY" },
    { label: "CT", value: "CT" },
    { label: "IR", value: "IR" },
    { label: "MRI", value: "MRI" },
    { label: "US", value: "US" },
    { label: "Dexa", value: "Dexa" },
    { label: "Nuc Med", value: "Nuc Med" },
  ]);

  useEffect(() => {
    (async () => {
      const { status } = await BarCodeScanner.requestPermissionsAsync();
      setHasPermission(status === "granted");
    })();
  }, []);

  const handleBarCodeScanned = ({ type, data }) => {
    setScanned(true);
    // send message to server with scanned patient ID
    alert(`${data}`);
    fetch(`http://localhost:8080/scan-exams?modality=${modality}&patientID=${data}`);
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
        onBarCodeScanned={scanned || !modality ? undefined : handleBarCodeScanned}
        style={StyleSheet.absoluteFillObject}
      />
      <View style={styles.selectModality}>
        {scanned && (
          <Button title={"Tap to Scan Again"} onPress={() => setScanned(false)} />
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
