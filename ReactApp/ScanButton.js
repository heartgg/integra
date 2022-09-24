import {Camera} from 'react-native-vision-camera'
async function hasCameraPerms() {
    const cameraPermission = await Camera.getCameraPermissionStatus()
    if (cameraPermission) {
        return true;
    }
    else {
        return await Camera.requestCameraPermission();
    }
}

function tryOpenCamera() {
    print("Returned from hasCameraPerms: ",hasCameraPerms());
    // if (!hasCameraPerms()) {return;};
  
}

export {hasCameraPerms, tryOpenCamera};