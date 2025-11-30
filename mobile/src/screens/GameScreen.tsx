import React, { useEffect, useState } from 'react';
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  ActivityIndicator,
  Alert,
} from 'react-native';
import { useGameStore } from '../store/gameStore';
import BLEService from '../services/BLEService';
import APIService from '../services/APIService';

export default function GameScreen() {
  const {
    isConnected,
    deviceName,
    deviceId,
    gameStatus,
    timerValue,
    batteryLevel,
    setConnected,
    setGameStatus,
    setTimerValue,
    setBatteryLevel,
  } = useGameStore();

  const [scanning, setScanning] = useState(false);
  const [connecting, setConnecting] = useState(false);

  useEffect(() => {
    initializeBLE();

    return () => {
      if (isConnected) {
        BLEService.disconnect();
      }
    };
  }, []);

  useEffect(() => {
    if (isConnected) {
      subscribeToUpdates();
    }
  }, [isConnected]);

  const initializeBLE = async () => {
    const initialized = await BLEService.initialize();
    if (!initialized) {
      Alert.alert('Error', 'Please enable Bluetooth');
    }
  };

  const startScan = async () => {
    setScanning(true);

    await BLEService.scanDevices((device) => {
      setScanning(false);

      Alert.alert(
        'Device Found',
        `Connect to ${device.name}?`,
        [
          { text: 'Cancel', style: 'cancel' },
          {
            text: 'Connect',
            onPress: () => connectToDevice(device.id),
          },
        ]
      );
    }, 10000);

    setTimeout(() => setScanning(false), 10000);
  };

  const connectToDevice = async (deviceId: string) => {
    setConnecting(true);

    const connected = await BLEService.connect(deviceId);

    if (connected) {
      const name = BLEService.getDeviceName();
      setConnected(true, name || undefined);

      // Save device ID
      await APIService.saveDeviceId(deviceId);

      Alert.alert('Success', 'Connected to device!');
    } else {
      Alert.alert('Error', 'Failed to connect to device');
    }

    setConnecting(false);
  };

  const subscribeToUpdates = async () => {
    try {
      // Subscribe to timer updates
      await BLEService.subscribeToTimer((timer) => {
        setTimerValue(timer);
      });

      // Subscribe to status updates
      await BLEService.subscribeToStatus((status) => {
        console.log('Status update:', status);
        // Parse status and update game state
      });

      console.log('Subscribed to device updates');
    } catch (error) {
      console.error('Subscribe error:', error);
    }
  };

  const disconnect = async () => {
    await BLEService.disconnect();
    setConnected(false, undefined);
  };

  const formatTime = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  };

  const getStatusColor = () => {
    switch (gameStatus) {
      case 'ready':
        return '#4CAF50';
      case 'playing':
        return '#2196F3';
      case 'completed':
        return '#FF9800';
      default:
        return '#9E9E9E';
    }
  };

  return (
    <View style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Text style={styles.title}>Maze Challenge</Text>
        {isConnected && (
          <View style={styles.connectionInfo}>
            <View style={styles.connectedDot} />
            <Text style={styles.deviceName}>{deviceName}</Text>
          </View>
        )}
      </View>

      {/* Main Content */}
      {!isConnected ? (
        <View style={styles.disconnectedContainer}>
          <Text style={styles.disconnectedText}>
            Connect to your maze device
          </Text>

          <TouchableOpacity
            style={[styles.button, scanning && styles.buttonDisabled]}
            onPress={startScan}
            disabled={scanning || connecting}
          >
            {scanning ? (
              <ActivityIndicator color="#fff" />
            ) : (
              <Text style={styles.buttonText}>
                {connecting ? 'Connecting...' : 'Scan for Devices'}
              </Text>
            )}
          </TouchableOpacity>
        </View>
      ) : (
        <View style={styles.gameContainer}>
          {/* Status */}
          <View
            style={[
              styles.statusBadge,
              { backgroundColor: getStatusColor() },
            ]}
          >
            <Text style={styles.statusText}>
              {gameStatus.toUpperCase()}
            </Text>
          </View>

          {/* Timer */}
          <View style={styles.timerContainer}>
            <Text style={styles.timerLabel}>TIME</Text>
            <Text style={styles.timerValue}>{formatTime(timerValue)}</Text>
          </View>

          {/* Battery */}
          <View style={styles.batteryContainer}>
            <Text style={styles.batteryLabel}>Battery</Text>
            <View style={styles.batteryBar}>
              <View
                style={[
                  styles.batteryFill,
                  {
                    width: `${batteryLevel}%`,
                    backgroundColor:
                      batteryLevel > 20 ? '#4CAF50' : '#F44336',
                  },
                ]}
              />
            </View>
            <Text style={styles.batteryText}>{batteryLevel}%</Text>
          </View>

          {/* Instructions */}
          <View style={styles.instructionsContainer}>
            {gameStatus === 'ready' && (
              <Text style={styles.instructions}>
                Place the ball at the START sensor
              </Text>
            )}
            {gameStatus === 'playing' && (
              <Text style={styles.instructions}>
                Navigate through the maze!
              </Text>
            )}
            {gameStatus === 'completed' && (
              <Text style={styles.instructions}>
                Completed! Press button to reset
              </Text>
            )}
          </View>

          {/* Disconnect Button */}
          <TouchableOpacity
            style={styles.disconnectButton}
            onPress={disconnect}
          >
            <Text style={styles.disconnectButtonText}>Disconnect</Text>
          </TouchableOpacity>
        </View>
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  header: {
    backgroundColor: '#6200EA',
    padding: 20,
    paddingTop: 60,
    alignItems: 'center',
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    color: '#fff',
  },
  connectionInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    marginTop: 10,
  },
  connectedDot: {
    width: 8,
    height: 8,
    borderRadius: 4,
    backgroundColor: '#4CAF50',
    marginRight: 8,
  },
  deviceName: {
    color: '#fff',
    fontSize: 14,
  },
  disconnectedContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
  },
  disconnectedText: {
    fontSize: 18,
    color: '#666',
    marginBottom: 30,
    textAlign: 'center',
  },
  button: {
    backgroundColor: '#6200EA',
    paddingHorizontal: 40,
    paddingVertical: 15,
    borderRadius: 25,
    minWidth: 200,
    alignItems: 'center',
  },
  buttonDisabled: {
    backgroundColor: '#9E9E9E',
  },
  buttonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
  gameContainer: {
    flex: 1,
    padding: 20,
  },
  statusBadge: {
    alignSelf: 'center',
    paddingHorizontal: 20,
    paddingVertical: 10,
    borderRadius: 20,
    marginTop: 20,
  },
  statusText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: 'bold',
  },
  timerContainer: {
    alignItems: 'center',
    marginTop: 40,
  },
  timerLabel: {
    fontSize: 18,
    color: '#666',
    marginBottom: 10,
  },
  timerValue: {
    fontSize: 64,
    fontWeight: 'bold',
    color: '#333',
    fontVariant: ['tabular-nums'],
  },
  batteryContainer: {
    marginTop: 40,
    alignItems: 'center',
  },
  batteryLabel: {
    fontSize: 16,
    color: '#666',
    marginBottom: 10,
  },
  batteryBar: {
    width: 200,
    height: 20,
    backgroundColor: '#E0E0E0',
    borderRadius: 10,
    overflow: 'hidden',
  },
  batteryFill: {
    height: '100%',
  },
  batteryText: {
    marginTop: 5,
    fontSize: 14,
    color: '#666',
  },
  instructionsContainer: {
    marginTop: 40,
    padding: 20,
    backgroundColor: '#fff',
    borderRadius: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 3,
  },
  instructions: {
    fontSize: 16,
    color: '#333',
    textAlign: 'center',
  },
  disconnectButton: {
    marginTop: 'auto',
    backgroundColor: '#F44336',
    paddingVertical: 15,
    borderRadius: 25,
    alignItems: 'center',
  },
  disconnectButtonText: {
    color: '#fff',
    fontSize: 16,
    fontWeight: '600',
  },
});
