import { create } from 'zustand';

interface GameState {
  // Connection
  isConnected: boolean;
  deviceName: string | null;
  deviceId: string | null;

  // Game state
  gameStatus: 'idle' | 'ready' | 'playing' | 'completed';
  timerValue: number;
  batteryLevel: number;

  // Statistics
  totalGames: number;
  completedGames: number;
  averageTime: number;
  bestTime: number;

  // Actions
  setConnected: (connected: boolean, deviceName?: string) => void;
  setDeviceId: (deviceId: string) => void;
  setGameStatus: (status: 'idle' | 'ready' | 'playing' | 'completed') => void;
  setTimerValue: (value: number) => void;
  setBatteryLevel: (level: number) => void;
  setStatistics: (stats: {
    totalGames: number;
    completedGames: number;
    averageTime: number;
    bestTime: number;
  }) => void;
  reset: () => void;
}

export const useGameStore = create<GameState>((set) => ({
  // Initial state
  isConnected: false,
  deviceName: null,
  deviceId: null,
  gameStatus: 'idle',
  timerValue: 0,
  batteryLevel: 100,
  totalGames: 0,
  completedGames: 0,
  averageTime: 0,
  bestTime: 0,

  // Actions
  setConnected: (connected, deviceName) =>
    set({ isConnected: connected, deviceName: deviceName || null }),

  setDeviceId: (deviceId) => set({ deviceId }),

  setGameStatus: (status) => set({ gameStatus: status }),

  setTimerValue: (value) => set({ timerValue: value }),

  setBatteryLevel: (level) => set({ batteryLevel: level }),

  setStatistics: (stats) =>
    set({
      totalGames: stats.totalGames,
      completedGames: stats.completedGames,
      averageTime: stats.averageTime,
      bestTime: stats.bestTime,
    }),

  reset: () =>
    set({
      isConnected: false,
      deviceName: null,
      gameStatus: 'idle',
      timerValue: 0,
    }),
}));
