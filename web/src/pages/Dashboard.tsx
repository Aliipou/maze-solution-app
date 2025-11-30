import React, { useEffect, useState } from 'react';
import {
  Box,
  Grid,
  Paper,
  Typography,
  Card,
  CardContent,
  CircularProgress,
} from '@mui/material';
import {
  TrendingUp,
  Devices as DevicesIcon,
  Timer,
  CheckCircle,
} from '@mui/icons-material';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import apiService from '../services/apiService';

interface Stats {
  totalDevices: number;
  activeGames: number;
  averageTime: number;
  completedToday: number;
}

export default function Dashboard() {
  const [stats, setStats] = useState<Stats>({
    totalDevices: 0,
    activeGames: 0,
    averageTime: 0,
    completedToday: 0,
  });
  const [loading, setLoading] = useState(true);
  const [chartData, setChartData] = useState<any[]>([]);

  useEffect(() => {
    loadDashboardData();
    const interval = setInterval(loadDashboardData, 30000); // Refresh every 30s
    return () => clearInterval(interval);
  }, []);

  const loadDashboardData = async () => {
    try {
      const statuses = await apiService.getDeviceStatuses();

      // Calculate stats
      const devices = new Set(statuses.map(s => s.device_id));
      const activeGames = statuses.filter(s => s.alarm_active).length;
      const completedToday = statuses.filter(s =>
        s.maze_completed &&
        new Date(s.timestamp).toDateString() === new Date().toDateString()
      ).length;

      setStats({
        totalDevices: devices.size,
        activeGames,
        averageTime: 120, // Would calculate from actual data
        completedToday,
      });

      // Generate chart data (last 7 days)
      const last7Days = Array.from({ length: 7 }, (_, i) => {
        const date = new Date();
        date.setDate(date.getDate() - (6 - i));
        return {
          name: date.toLocaleDateString('en-US', { weekday: 'short' }),
          games: Math.floor(Math.random() * 20) + 5, // Mock data
          completed: Math.floor(Math.random() * 15) + 3,
        };
      });

      setChartData(last7Days);
      setLoading(false);
    } catch (error) {
      console.error('Load dashboard error:', error);
      setLoading(false);
    }
  };

  const StatCard = ({ title, value, icon: Icon, color }: any) => (
    <Card sx={{ height: '100%' }}>
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="flex-start">
          <Box>
            <Typography color="textSecondary" gutterBottom variant="body2">
              {title}
            </Typography>
            <Typography variant="h4" component="div">
              {value}
            </Typography>
          </Box>
          <Box
            sx={{
              backgroundColor: `${color}20`,
              borderRadius: 2,
              p: 1,
            }}
          >
            <Icon sx={{ color, fontSize: 32 }} />
          </Box>
        </Box>
      </CardContent>
    </Card>
  );

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="80vh">
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>

      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Total Devices"
            value={stats.totalDevices}
            icon={DevicesIcon}
            color="#6200EA"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Active Games"
            value={stats.activeGames}
            icon={TrendingUp}
            color="#2196F3"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Avg Time (sec)"
            value={stats.averageTime}
            icon={Timer}
            color="#FF9800"
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <StatCard
            title="Completed Today"
            value={stats.completedToday}
            icon={CheckCircle}
            color="#4CAF50"
          />
        </Grid>
      </Grid>

      {/* Charts */}
      <Grid container spacing={3}>
        {/* Line Chart */}
        <Grid item xs={12} lg={8}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Activity Over Last 7 Days
            </Typography>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Line type="monotone" dataKey="games" stroke="#6200EA" name="Total Games" />
                <Line type="monotone" dataKey="completed" stroke="#4CAF50" name="Completed" />
              </LineChart>
            </ResponsiveContainer>
          </Paper>
        </Grid>

        {/* Bar Chart */}
        <Grid item xs={12} lg={4}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Completion Rate
            </Typography>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis />
                <Tooltip />
                <Bar dataKey="completed" fill="#4CAF50" />
              </BarChart>
            </ResponsiveContainer>
          </Paper>
        </Grid>

        {/* Recent Activity */}
        <Grid item xs={12}>
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom>
              Recent Activity
            </Typography>
            <Box>
              {[1, 2, 3, 4, 5].map((i) => (
                <Box
                  key={i}
                  sx={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    py: 2,
                    borderBottom: i !== 5 ? '1px solid #eee' : 'none',
                  }}
                >
                  <Box>
                    <Typography variant="body1" fontWeight={500}>
                      Device ESP32_MAZE_00{i}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      Completed maze in 1:45
                    </Typography>
                  </Box>
                  <Typography variant="body2" color="textSecondary">
                    {i} min ago
                  </Typography>
                </Box>
              ))}
            </Box>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
}
