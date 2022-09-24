import { StatusBar } from 'expo-status-bar';
import { StyleSheet, Text, View } from 'react-native';
import DropDownPicker from 'react-native-dropdown-picker';
import { useState } from 'react';
import ScanButton from './ScanButton.js';

export default function App() {
  const [open, setOpen] = useState(false);
  const [value, setValue] = useState(null);
  const [items, setItems] = useState([
    {label: 'IE Fluoro', value: 'IE Fluoro'},
    {label: 'XRAY', value: 'XRAY'},
    {label: 'CT', value: 'CT'},
    {label: 'IR', value: 'IR'},
    {label: 'MRI', value: 'MRI'},
    {label: 'US', value: 'US'},
    {label: 'Dexa', value: 'Dexa'},
    {label: 'Nuc Med', value: 'Nuc Med'}
  ]);

  return (
    <View style={styles.container}>
      <DropDownPicker
        open={open}
        value={value}
        items={items}
        setOpen={setOpen}
        setValue={setValue}
        setItems={setItems}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
