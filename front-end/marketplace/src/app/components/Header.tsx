"use client";

import React, { useContext, useEffect } from 'react';
import Tabs from '@mui/material/Tabs';
import Tab from '@mui/material/Tab';
import Button from '@mui/material/Button';
import AppBar from '@mui/material/AppBar';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import MenuItem from '@mui/material/MenuItem';
import Menu from '@mui/material/Menu';
import IconButton from '@mui/material/IconButton';
import AccountCircle from '@mui/icons-material/AccountCircle';
import Badge from '@mui/material/Badge';
import NotificationsIcon from '@mui/icons-material/Notifications';
import Popover from '@mui/material/Popover';
import Avatar from '@mui/material/Avatar';
import List from '@mui/material/List';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import ListItemText from '@mui/material/ListItemText';
import ListItem from '@mui/material/ListItem';
import ImageIcon from '@mui/icons-material/Image';
import dayjs from 'dayjs';
import Visibility from '@mui/icons-material/Visibility';
import Delete from '@mui/icons-material/Delete';
import { SnackbarContext, SnackbarProviderContextType } from '../context/SnackbarContext';
import { SocketContext } from '../context/SocketContext';
import { useFetchUser, useLogoutUser } from '../requests/userHooks';
import { useDeleteNotification, useFetchNotifications, useWatchNotification } from '../requests/notificationHooks';

function Header() {
  const { openSuccessSnackbar, openErrorSnackbar }: Partial<SnackbarProviderContextType> = useContext(SnackbarContext);
  const { socket, setUserStatus } = useContext(SocketContext);
  const { isLoadingUser, user, removeUserData } = useFetchUser();
  const { refetchNotifications, notifications } = useFetchNotifications({ user });
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const [anchorElNotifications, setAnchorElNotifications] = React.useState<null | HTMLElement>(null);

  useEffect(() => {
    console.log("🚀 ~ file: Header.tsx:41 ~ useEffect ~ setUserStatus:", setUserStatus)
    if (setUserStatus) setUserStatus(Boolean(user));
  }, [user, setUserStatus]);

  useEffect(() => {
    if (socket && user) {
      socket.on('notifications_update', () => {
        refetchNotifications();
      });
    }
    if (socket && !user) {
      socket.off('notifications_update');
    }
    return () => {
      if (socket) socket.off('notifications_update');
    };
  }, [socket, user, openSuccessSnackbar, refetchNotifications]);

  const isLoggedIn = Boolean(user);
  const menuId = 'menu';
  const tabs = [
    { path: '/', label: 'Home' },
    { path: '/chat', label: 'Chat' },
  ];
  const router = useRouter();
  const [tabIndex, setTabIndex] = React.useState(false);
  const { pathname } = router;
  useEffect(() => {
    let activeTabIndex = null;
    const activeTab = tabs.find((tab, index) => {
      const isMatched = tab.path === pathname;
      if (isMatched) activeTabIndex = index;
      return tab.path === pathname;
    });

    if (!activeTab) {
      setTabIndex(false);
    } else {
      setTabIndex(activeTabIndex);
    }
  }, [tabs, pathname]);
  const { logoutUser } = useLogoutUser({
    onSuccess: () => {
      removeUserData();
      openSuccessSnackbar('Logged out');
    },
    onError: (error) => {
      openErrorSnackbar(error?.response?.data?.message || error?.message);
    },
  });

  const { viewNotification } = useWatchNotification({
    onError: (error) => {
      openErrorSnackbar(error?.response?.data?.message || error?.message);
    },
  });

  const { deleteNotification } = useDeleteNotification({
    onError: (error) => {
      openErrorSnackbar(error?.response?.data?.message || error?.message);
    },
  });

  const handleChange = (event: React.MouseEvent<HTMLButtonElement>, newValue: number) => {
    const tab = tabs[newValue];
    const { path } = tab;
    router.push(path);
  };

  const handleClickNotifications = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorElNotifications(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorElNotifications(null);
  };
  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const isMenuOpen = Boolean(anchorEl);
  const open = Boolean(anchorElNotifications);
  const id = open ? 'simple-popover' : undefined;

  const handleOpenMyAccount = () => {
    router.push('/account');
  };

  const renderMenu = (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{
        vertical: 'top',
        horizontal: 'right',
      }}
      id={menuId}
      keepMounted
      transformOrigin={{
        vertical: 'top',
        horizontal: 'right',
      }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem
        onClick={() => {
          handleMenuClose();
          handleOpenMyAccount();
        }}
      >
        My Account
      </MenuItem>
      <MenuItem
        onClick={() => {
          handleMenuClose();
          logoutUser();
        }}
      >
        Logout
      </MenuItem>
    </Menu>
  );

  const notViewedNotificationsNumber = 0 || notifications?.filter((n) => n.viewedAt === '0001-01-01T00:00:00Z').length;

  const renderPopover = (
    <Popover
      id={id}
      open={open}
      anchorEl={anchorElNotifications}
      onClose={handleClose}
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'left',
      }}
    >
      {/* deleteNotification */}
      <List sx={{ width: '100%', maxWidth: 360, bgcolor: 'background.paper' }}>
        {notifications?.map((n) => (
          <ListItem
            key={n.id}
            secondaryAction={
              n.viewedAt === '0001-01-01T00:00:00Z' ? (
                <IconButton size="large" edge="end" aria-label="view" onClick={() => viewNotification({ notificationId: n.id })}>
                  <Visibility className="text-24" />
                </IconButton>
              ) : (
                <IconButton
                  size="large"
                  edge="end"
                  aria-label="delete"
                  onClick={() => deleteNotification({ notificationId: n.id })}
                >
                  <Delete className="text-24" />
                </IconButton>
              )
            }
          >
            <ListItemAvatar>
              <Avatar>
                <ImageIcon />
              </Avatar>
            </ListItemAvatar>
            <ListItemText primary={`${n.title}  ${dayjs(n.createdAt).format('HH:mm:ss MMMM DD, YYYY')}`} secondary={n.message} />
          </ListItem>
        ))}
      </List>
    </Popover>
  );

  return (
    <>
      <AppBar color="default" className="h-48 px-0 sm:px-16">
        <div className="w-full flex justify-between items-center">
          <Tabs
            value={tabIndex}
            indicatorColor="secondary"
            textColor="secondary"
            onChange={handleChange}
            aria-label="disabled tabs example"
          >
            {tabs.map(({ label }) => (
              <Tab key={label} label={label} />
            ))}
          </Tabs>
          {!isLoadingUser && (
            <div className="mr-8 sm:mr-0">
              {!isLoggedIn ? (
                <div className="flex ml-auto text-8 sm:text-14">
                  {pathname !== '/signin' && (
                    <div className="ml-16">
                      <Link href="/signin" className="no-underline">
                        <Button
                          variant="outlined"
                          onClick={() => {
                            router.push('/signin');
                          }}
                        >
                          Sign In
                        </Button>
                      </Link>
                    </div>
                  )}
                  {pathname !== '/signup' && (
                    <div className="ml-8">
                      <Link href="/signup" className="no-underline">
                        <Button
                          variant="outlined"
                          onClick={() => {
                            router.push('/signup');
                          }}
                        >
                          Sign Up
                        </Button>
                      </Link>
                    </div>
                  )}
                </div>
              ) : (
                <div className="flex items-center ml-auto">
                  <div className="text-14 hidden sm:block mr-16">{user?.name}</div>
                  <IconButton
                    className="mr-16"
                    size="large"
                    aria-label="notifications"
                    color="primary"
                    onClick={handleClickNotifications}
                  >
                    <Badge badgeContent={notViewedNotificationsNumber} color="error">
                      <NotificationsIcon className="text-24" />
                    </Badge>
                  </IconButton>
                  <IconButton
                    size="small"
                    aria-label="account of current user"
                    aria-controls={menuId}
                    aria-haspopup="true"
                    onClick={handleProfileMenuOpen}
                    color="primary"
                  >
                    <AccountCircle className="text-24" />
                  </IconButton>
                </div>
              )}
            </div>
          )}
        </div>
      </AppBar>
      {renderMenu}
      {renderPopover}
      <div className="h-48" />
    </>
  );
}

export default Header;
