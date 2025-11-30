import { useState, useEffect } from 'react';
import {
  Container,
  Grid,
  Card,
  CardContent,
  Typography,
  Button,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
  Box,
  AppBar,
  Toolbar,
} from '@mui/material';
import {
  Refresh as RefreshIcon,
  PlayArrow as PlayIcon,
  Add as AddIcon,
  Delete as DeleteIcon,
} from '@mui/icons-material';
import axios from 'axios';

const API_URL = 'http://localhost:8080';
const AUTH = 'Basic ' + btoa('admin:password');

interface DeviceStatus {
  id: number;
  device_id: string;
  alarm_active: boolean;
  maze_completed: boolean;
  hall_sensor_value: boolean;
  battery_level: number;
  timestamp: string;
}

function App() {
  const [sessions, setSessions] = useState<DeviceStatus[]>([]);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState({
    totalDevices: 0,
    activeGames: 0,
    completedToday: 0,
    avgBattery: 0,
  });

  const fetchData = async () => {
    setLoading(true);
    try {
      const response = await axios.get(`${API_URL}/device/status`, {
        headers: { Authorization: AUTH },
      });
      const data = response.data || [];
      setSessions(data);

      const devices = new Set(data.map((s: DeviceStatus) => s.device_id));
      const activeGames = data.filter((s: DeviceStatus) => s.alarm_active).length;
      const today = new Date().toISOString().split('T')[0];
      const completedToday = data.filter(
        (s: DeviceStatus) => s.maze_completed && s.timestamp.startsWith(today)
      ).length;
      const avgBattery = data.length > 0
        ? Math.round(data.reduce((sum: number, s: DeviceStatus) => sum + s.battery_level, 0) / data.length)
        : 0;

      setStats({
        totalDevices: devices.size,
        activeGames,
        completedToday,
        avgBattery,
      });
    } catch (error) {
      console.error('Error fetching data:', error);
    } finally {
      setLoading(false);
    }
  };

  const createTestSession = async () => {
    try {
      const session = {
        device_id: 'DEMO_' + Math.random().toString(36).substr(2, 9).toUpperCase(),
        alarm_active: false,
        maze_completed: false,
        hall_sensor_value: false,
        battery_level: Math.floor(Math.random() * 20) + 80,
        timestamp: new Date().toISOString(),
      };
      await axios.post(`${API_URL}/device/status`, session, {
        headers: { Authorization: AUTH, 'Content-Type': 'application/json' },
      });
      fetchData();
    } catch (error) {
      console.error('Error creating session:', error);
    }
  };

  const simulateGame = async () => {
    try {
      const deviceId = 'GAME_SIM_' + Date.now();
      await axios.post(
        `${API_URL}/device/status`,
        {
          device_id: deviceId,
          alarm_active: true,
          maze_completed: false,
          hall_sensor_value: true,
          battery_level: 100,
          timestamp: new Date().toISOString(),
        },
        { headers: { Authorization: AUTH, 'Content-Type': 'application/json' } }
      );
      fetchData();
      setTimeout(async () => {
        await axios.post(
          `${API_URL}/device/status`,
          {
            device_id: deviceId,
            alarm_active: false,
            maze_completed: true,
            hall_sensor_value: true,
            battery_level: 98,
            timestamp: new Date().toISOString(),
          },
          { headers: { Authorization: AUTH, 'Content-Type': 'application/json' } }
        );
        fetchData();
      }, 2000);
    } catch (error) {
      console.error('Error simulating game:', error);
    }
  };

  const clearAll = async () => {
    if (!window.confirm('Clear all sessions?')) return;
    try {
      const response = await axios.get(`${API_URL}/device/status`, {
        headers: { Authorization: AUTH },
      });
      const data = response.data || [];
      for (const session of data) {
        await axios.delete(`${API_URL}/device/status/${session.id}`, {
          headers: { Authorization: AUTH },
        });
      }
      fetchData();
    } catch (error) {
      console.error('Error clearing sessions:', error);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            🎮 Maze Challenge Dashboard
          </Typography>
          <Chip label="LIVE" color="success" size="small" />
        </Toolbar>
      </AppBar>

      <Container maxWidth="xl" sx={{ mt: 4, mb: 4 }}>
        <Grid container spacing={3} sx={{ mb: 3 }}>
          <Grid item xs={12} sm={6} md={3}>
            <Card><CardContent>
                <Typography color="textSecondary" gutterBottom>Total Devices</Typography>
                <Typography variant="h3">{stats.totalDevices}</Typography>
              </CardContent></Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card><CardContent>
                <Typography color="textSecondary" gutterBottom>Active Games</Typography>
                <Typography variant="h3">{stats.activeGames}</Typography>
              </CardContent></Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card><CardContent>
                <Typography color="textSecondary" gutterBottom>Completed Today</Typography>
                <Typography variant="h3">{stats.completedToday}</Typography>
              </CardContent></Card>
          </Grid>
          <Grid item xs={12} sm={6} md={3}>
            <Card><CardContent>
                <Typography color="textSecondary" gutterBottom>Avg Battery</Typography>
                <Typography variant="h3">{stats.avgBattery}%</Typography>
              </CardContent></Card>
          </Grid>
        </Grid>

        <Card sx={{ mb: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>API Controls</Typography>
            <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
              <Button variant="contained" startIcon={<RefreshIcon />} onClick={fetchData} disabled={loading}>Refresh</Button>
              <Button variant="contained" startIcon={<AddIcon />} onClick={createTestSession} color="secondary">Create Test Session</Button>
              <Button variant="contained" startIcon={<PlayIcon />} onClick={simulateGame} color="success">Simulate Game</Button>
              <Button variant="outlined" startIcon={<DeleteIcon />} onClick={clearAll} color="error">Clear All</Button>
            </Box>
          </CardContent>
        </Card>

        <Card>
          <CardContent>
            <Typography variant="h6" gutterBottom>Recent Game Sessions</Typography>
            <TableContainer component={Paper}>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>Device ID</TableCell>
                    <TableCell>Status</TableCell>
                    <TableCell>Battery</TableCell>
                    <TableCell>Sensor</TableCell>
                    <TableCell>Timestamp</TableCell>
                  </TableRow>
                </TableHead>
                <TableBody>
                  {sessions.slice(0, 10).map((session) => (
                    <TableRow key={session.id}>
                      <TableCell><strong>{session.device_id}</strong></TableCell>
                      <TableCell>
                        {session.maze_completed && <Chip label="Completed" color="success" size="small" />}
                        {session.alarm_active && !session.maze_completed && <Chip label="Active" color="info" size="small" />}
                        {!session.alarm_active && !session.maze_completed && <Chip label="Idle" color="warning" size="small" />}
                      </TableCell>
                      <TableCell>{session.battery_level}%</TableCell>
                      <TableCell>{session.hall_sensor_value ? '✓' : '✗'}</TableCell>
                      <TableCell>{new Date(session.timestamp).toLocaleString()}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
            {sessions.length === 0 && (
              <Typography variant="body2" color="textSecondary" align="center" sx={{ mt: 2 }}>
                No sessions yet. Create one to get started!
              </Typography>
            )}
          </CardContent>
        </Card>

        <Card sx={{ mt: 3 }}>
          <CardContent>
            <Typography variant="h6" gutterBottom>System Information</Typography>
            <Typography variant="body2"><strong>Backend API:</strong> <Chip label="Running" color="success" size="small" /> - http://localhost:8080</Typography>
            <Typography variant="body2"><strong>Database:</strong> SQLite (production.db)</Typography>
            <Typography variant="body2"><strong>Authentication:</strong> Basic Auth (admin)</Typography>
            <Typography variant="body2"><strong>Last Updated:</strong> {new Date().toLocaleTimeString()}</Typography>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
}

export default App;
