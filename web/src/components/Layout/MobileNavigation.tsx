import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  BottomNavigation,
  BottomNavigationAction,
  Badge,
  Fab,
  SpeedDial,
  SpeedDialAction,
  SpeedDialIcon,
  useMediaQuery,
  useTheme
} from '@mui/material';
import {
  Home as HomeIcon,
  CalendarMonth as CalendarIcon,
  Flag as GoalsIcon,
  Mood as MoodIcon,
  Settings as SettingsIcon,
  Add as AddIcon,
  Event as EventIcon,
  EmojiEvents as GoalIcon,
  SentimentSatisfied as RecordMoodIcon
} from '@mui/icons-material';
import { Paper, Box } from '@mui/material';

interface MobileNavigationProps {
  unreadNotifications?: number;
  todayEvents?: number;
  pendingGoals?: number;
}

const MobileNavigation: React.FC<MobileNavigationProps> = ({
  unreadNotifications = 0,
  todayEvents = 0,
  pendingGoals = 0
}) => {
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  
  const [speedDialOpen, setSpeedDialOpen] = useState(false);

  // Don't render on desktop
  if (!isMobile) return null;

  const getCurrentTab = () => {
    const path = location.pathname;
    if (path === '/' || path === '/dashboard') return 0;
    if (path === '/calendar') return 1;
    if (path === '/goals') return 2;
    if (path === '/moods') return 3;
    if (path === '/settings') return 4;
    return 0;
  };

  const handleNavigationChange = (event: React.SyntheticEvent, newValue: number) => {
    const routes = ['/', '/calendar', '/goals', '/moods', '/settings'];
    navigate(routes[newValue]);
  };

  const speedDialActions = [
    {
      icon: <EventIcon />,
      name: 'New Event',
      onClick: () => {
        navigate('/calendar');
        setSpeedDialOpen(false);
        // Trigger event creation modal
        setTimeout(() => {
          const createEventBtn = document.querySelector('[data-testid="create-event-btn"]') as HTMLElement;
          createEventBtn?.click();
        }, 100);
      }
    },
    {
      icon: <GoalIcon />,
      name: 'New Goal',
      onClick: () => {
        navigate('/goals');
        setSpeedDialOpen(false);
        // Trigger goal creation modal
        setTimeout(() => {
          const createGoalBtn = document.querySelector('[data-testid="create-goal-btn"]') as HTMLElement;
          createGoalBtn?.click();
        }, 100);
      }
    },
    {
      icon: <RecordMoodIcon />,
      name: 'Record Mood',
      onClick: () => {
        navigate('/moods');
        setSpeedDialOpen(false);
        // Trigger mood recording modal
        setTimeout(() => {
          const recordMoodBtn = document.querySelector('[data-testid="record-mood-btn"]') as HTMLElement;
          recordMoodBtn?.click();
        }, 100);
      }
    }
  ];

  return (
    <>
      {/* Bottom Navigation Bar */}
      <Paper 
        sx={{ 
          position: 'fixed', 
          bottom: 0, 
          left: 0, 
          right: 0, 
          zIndex: 1000,
          borderTop: '1px solid #f0f0f0',
          boxShadow: '0 -2px 8px rgba(0, 0, 0, 0.1)'
        }} 
        elevation={3}
      >
        <BottomNavigation
          value={getCurrentTab()}
          onChange={handleNavigationChange}
          showLabels
          sx={{
            height: 64,
            '& .MuiBottomNavigationAction-root': {
              minWidth: 'auto',
              padding: '6px 12px 8px',
              '&.Mui-selected': {
                color: '#1677ff'
              }
            }
          }}
        >
          <BottomNavigationAction 
            label="Home" 
            icon={<HomeIcon />} 
          />
          <BottomNavigationAction 
            label="Calendar" 
            icon={
              <Badge badgeContent={todayEvents > 0 ? todayEvents : null} color="primary">
                <CalendarIcon />
              </Badge>
            } 
          />
          <BottomNavigationAction 
            label="Goals" 
            icon={
              <Badge badgeContent={pendingGoals > 0 ? pendingGoals : null} color="secondary">
                <GoalsIcon />
              </Badge>
            } 
          />
          <BottomNavigationAction 
            label="Mood" 
            icon={<MoodIcon />} 
          />
          <BottomNavigationAction 
            label="Settings" 
            icon={
              <Badge badgeContent={unreadNotifications > 0 ? unreadNotifications : null} color="error">
                <SettingsIcon />
              </Badge>
            } 
          />
        </BottomNavigation>
      </Paper>

      {/* Floating Action Button with Speed Dial */}
      <Box
        sx={{
          position: 'fixed',
          bottom: 80, // Above bottom navigation
          right: 16,
          zIndex: 1001
        }}
      >
        <SpeedDial
          ariaLabel="Quick Actions"
          sx={{
            '& .MuiFab-primary': {
              backgroundColor: '#1677ff',
              '&:hover': {
                backgroundColor: '#0958d9'
              }
            }
          }}
          icon={<SpeedDialIcon />}
          onClose={() => setSpeedDialOpen(false)}
          onOpen={() => setSpeedDialOpen(true)}
          open={speedDialOpen}
          direction="up"
        >
          {speedDialActions.map((action) => (
            <SpeedDialAction
              key={action.name}
              icon={action.icon}
              tooltipTitle={action.name}
              onClick={action.onClick}
              sx={{
                '& .MuiFab-primary': {
                  backgroundColor: '#52c41a',
                  '&:hover': {
                    backgroundColor: '#389e0d'
                  }
                }
              }}
            />
          ))}
        </SpeedDial>
      </Box>

      {/* Spacer to prevent content from being hidden behind bottom navigation */}
      <Box sx={{ height: 64 }} />
    </>
  );
};

export default MobileNavigation;