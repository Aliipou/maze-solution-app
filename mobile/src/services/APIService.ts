import axios, { AxiosInstance } from 'axios';
import AsyncStorage from '@react-native-async-storage/async-storage';

const API_BASE_URL = 'http://192.168.1.100:8080';
const USERNAME = 'admin';
const PASSWORD = 'password';

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

interface GameSession {
  id: number;
  device_id: string;
  start_time: string;
  end_time?: string;
  duration_ms?: number;
  completed: boolean;
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

  // ==================== Device Status ====================

  async getDeviceStatus(deviceId?: string): Promise<DeviceStatus[]> {
    try {
      const params = deviceId ? { device_id: deviceId } : {};
      const response = await this.client.get('/device/status', { params });
      return response.data;
    } catch (error) {
      console.error('Get device status error:', error);
      throw error;
    }
  }

  async createDeviceStatus(status: DeviceStatus): Promise<DeviceStatus> {
    try {
      const response = await this.client.post('/device/status', status);
      return response.data;
    } catch (error) {
      console.error('Create device status error:', error);
      throw error;
    }
  }

  async getLatestStatus(deviceId: string): Promise<DeviceStatus | null> {
    try {
      const statuses = await this.getDeviceStatus(deviceId);
      return statuses.length > 0 ? statuses[0] : null;
    } catch (error) {
      return null;
    }
  }

  // ==================== Device Config ====================

  async getDeviceConfig(deviceId: string): Promise<DeviceConfig | null> {
    try {
      const response = await this.client.get('/device/config', {
        params: { device_id: deviceId },
      });
      const configs = response.data;
      return configs.length > 0 ? configs[0] : null;
    } catch (error) {
      console.error('Get device config error:', error);
      return null;
    }
  }

  async updateDeviceConfig(config: DeviceConfig): Promise<DeviceConfig> {
    try {
      const response = await this.client.put('/device/config', config);
      return response.data;
    } catch (error) {
      console.error('Update device config error:', error);
      throw error;
    }
  }

  async createDeviceConfig(config: DeviceConfig): Promise<DeviceConfig> {
    try {
      const response = await this.client.post('/device/config', config);
      return response.data;
    } catch (error) {
      console.error('Create device config error:', error);
      throw error;
    }
  }

  // ==================== Statistics ====================

  async getGameSessions(deviceId: string, limit: number = 10): Promise<GameSession[]> {
    try {
      // This would be a new endpoint in the API
      // For now, we'll use device status as proxy
      const statuses = await this.getDeviceStatus(deviceId);

      const sessions: GameSession[] = statuses
        .filter(s => s.maze_completed)
        .slice(0, limit)
        .map((s, idx) => ({
          id: s.id || idx,
          device_id: s.device_id,
          start_time: s.timestamp,
          end_time: s.timestamp,
          duration_ms: 0, // Would need to calculate from data
          completed: s.maze_completed,
        }));

      return sessions;
    } catch (error) {
      console.error('Get game sessions error:', error);
      return [];
    }
  }

  async getStatistics(deviceId: string): Promise<{
    totalGames: number;
    completedGames: number;
    averageTime: number;
    bestTime: number;
  }> {
    try {
      const statuses = await this.getDeviceStatus(deviceId);

      const completedGames = statuses.filter(s => s.maze_completed);
      const totalGames = statuses.length;

      // Mock statistics (would need proper calculation)
      return {
        totalGames,
        completedGames: completedGames.length,
        averageTime: 120, // seconds
        bestTime: 45, // seconds
      };
    } catch (error) {
      console.error('Get statistics error:', error);
      return {
        totalGames: 0,
        completedGames: 0,
        averageTime: 0,
        bestTime: 0,
      };
    }
  }

  // ==================== Local Storage ====================

  async saveDeviceId(deviceId: string): Promise<void> {
    await AsyncStorage.setItem('device_id', deviceId);
  }

  async getDeviceId(): Promise<string | null> {
    return await AsyncStorage.getItem('device_id');
  }

  async saveApiConfig(baseUrl: string, username: string, password: string): Promise<void> {
    await AsyncStorage.setItem('api_base_url', baseUrl);
    await AsyncStorage.setItem('api_username', username);
    await AsyncStorage.setItem('api_password', password);
  }
}

export default new APIService();
export type { DeviceStatus, DeviceConfig, GameSession };
