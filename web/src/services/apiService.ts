import axios, { AxiosInstance } from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
const USERNAME = process.env.REACT_APP_API_USERNAME || 'admin';
const PASSWORD = process.env.REACT_APP_API_PASSWORD || 'password';

interface DeviceStatus {
  id?: number;
  device_id: string;
  alarm_active: boolean;
  maze_completed: boolean;
  hall_sensor_value: boolean;
  battery_level: number;
  timestamp: string;
}

interface DeviceConfig {
  id?: number;
  device_id: string;
  alarm_timeout: number;
  sensitivity_level: number;
  updated_at: string;
}

class APIService {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
      auth: {
        username: USERNAME,
        password: PASSWORD,
      },
    });
  }

  async getDeviceStatuses(deviceId?: string): Promise<DeviceStatus[]> {
    const params = deviceId ? { device_id: deviceId } : {};
    const response = await this.client.get('/device/status', { params });
    return response.data || [];
  }

  async getDeviceStatus(id: number): Promise<DeviceStatus> {
    const response = await this.client.get(`/device/status/${id}`);
    return response.data;
  }

  async createDeviceStatus(status: DeviceStatus): Promise<DeviceStatus> {
    const response = await this.client.post('/device/status', status);
    return response.data;
  }

  async updateDeviceStatus(status: DeviceStatus): Promise<DeviceStatus> {
    const response = await this.client.put('/device/status', status);
    return response.data;
  }

  async deleteDeviceStatus(id: number): Promise<void> {
    await this.client.delete(`/device/status/${id}`);
  }

  async getDeviceConfigs(deviceId?: string): Promise<DeviceConfig[]> {
    const params = deviceId ? { device_id: deviceId } : {};
    const response = await this.client.get('/device/config', { params });
    return response.data || [];
  }

  async getDeviceConfig(id: number): Promise<DeviceConfig> {
    const response = await this.client.get(`/device/config/${id}`);
    return response.data;
  }

  async createDeviceConfig(config: DeviceConfig): Promise<DeviceConfig> {
    const response = await this.client.post('/device/config', config);
    return response.data;
  }

  async updateDeviceConfig(config: DeviceConfig): Promise<DeviceConfig> {
    const response = await this.client.put('/device/config', config);
    return response.data;
  }

  async deleteDeviceConfig(id: number): Promise<void> {
    await this.client.delete(`/device/config/${id}`);
  }
}

export default new APIService();
export type { DeviceStatus, DeviceConfig };
