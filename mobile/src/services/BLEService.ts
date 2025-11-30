import { BleManager, Device, Characteristic } from 'react-native-ble-plx';

const SERVICE_UUID = '4fafc201-1fb5-459e-8fcc-c5c9c331914b';
const CHAR_STATUS_UUID = 'beb5483e-36e1-4688-b7f5-ea07361b26a8';
const CHAR_TIMER_UUID = 'beb5483e-36e1-4688-b7f5-ea07361b26a9';
const CHAR_CONTROL_UUID = 'beb5483e-36e1-4688-b7f5-ea07361b26aa';

class BLEService {
  private manager: BleManager;
  private connectedDevice: Device | null = null;

  constructor() {
    this.manager = new BleManager();
  }

  /**
   * Initialize BLE and request permissions
   */
  async initialize(): Promise<boolean> {
    try {
      const state = await this.manager.state();
      console.log('BLE State:', state);

      if (state !== 'PoweredOn') {
        console.warn('Bluetooth is not powered on');
        return false;
      }

      return true;
    } catch (error) {
      console.error('BLE initialization error:', error);
      return false;
    }
  }

  /**
   * Scan for maze devices
   */
  async scanDevices(
    onDeviceFound: (device: Device) => void,
    durationMs: number = 10000
  ): Promise<void> {
    console.log('Starting BLE scan...');

    this.manager.startDeviceScan(
      [SERVICE_UUID],
      { allowDuplicates: false },
      (error, device) => {
        if (error) {
          console.error('Scan error:', error);
          return;
        }

        if (device && device.name?.includes('MazeChallenge')) {
          console.log('Found device:', device.name, device.id);
          onDeviceFound(device);
        }
      }
    );

    // Stop scan after duration
    setTimeout(() => {
      this.manager.stopDeviceScan();
      console.log('Scan stopped');
    }, durationMs);
  }

  /**
   * Connect to a device
   */
  async connect(deviceId: string): Promise<boolean> {
    try {
      console.log('Connecting to device:', deviceId);

      const device = await this.manager.connectToDevice(deviceId, {
        timeout: 10000,
      });

      await device.discoverAllServicesAndCharacteristics();

      this.connectedDevice = device;
      console.log('Connected successfully');

      return true;
    } catch (error) {
      console.error('Connection error:', error);
      return false;
    }
  }

  /**
   * Disconnect from device
   */
  async disconnect(): Promise<void> {
    if (this.connectedDevice) {
      await this.connectedDevice.cancelConnection();
      this.connectedDevice = null;
      console.log('Disconnected');
    }
  }

  /**
   * Subscribe to timer updates
   */
  async subscribeToTimer(
    callback: (timerValue: number) => void
  ): Promise<void> {
    if (!this.connectedDevice) {
      throw new Error('No device connected');
    }

    this.connectedDevice.monitorCharacteristicForService(
      SERVICE_UUID,
      CHAR_TIMER_UUID,
      (error, characteristic) => {
        if (error) {
          console.error('Timer subscription error:', error);
          return;
        }

        if (characteristic?.value) {
          const timerValue = this.base64ToNumber(characteristic.value);
          callback(timerValue);
        }
      }
    );

    console.log('Subscribed to timer updates');
  }

  /**
   * Subscribe to status updates
   */
  async subscribeToStatus(
    callback: (status: string) => void
  ): Promise<void> {
    if (!this.connectedDevice) {
      throw new Error('No device connected');
    }

    this.connectedDevice.monitorCharacteristicForService(
      SERVICE_UUID,
      CHAR_STATUS_UUID,
      (error, characteristic) => {
        if (error) {
          console.error('Status subscription error:', error);
          return;
        }

        if (characteristic?.value) {
          const status = this.base64ToString(characteristic.value);
          callback(status);
        }
      }
    );

    console.log('Subscribed to status updates');
  }

  /**
   * Send control command
   */
  async sendCommand(command: string): Promise<void> {
    if (!this.connectedDevice) {
      throw new Error('No device connected');
    }

    const base64Value = this.stringToBase64(command);

    await this.connectedDevice.writeCharacteristicWithResponseForService(
      SERVICE_UUID,
      CHAR_CONTROL_UUID,
      base64Value
    );

    console.log('Command sent:', command);
  }

  /**
   * Check if device is connected
   */
  isConnected(): boolean {
    return this.connectedDevice !== null;
  }

  /**
   * Get connected device name
   */
  getDeviceName(): string | null {
    return this.connectedDevice?.name || null;
  }

  // ==================== Helper Methods ====================

  private base64ToString(base64: string): string {
    return Buffer.from(base64, 'base64').toString('utf-8');
  }

  private stringToBase64(str: string): string {
    return Buffer.from(str, 'utf-8').toString('base64');
  }

  private base64ToNumber(base64: string): number {
    const str = this.base64ToString(base64);
    return parseInt(str, 10) || 0;
  }
}

export default new BLEService();
