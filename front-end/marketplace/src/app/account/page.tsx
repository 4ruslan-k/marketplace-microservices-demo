"use client";

import React, { useContext, useEffect } from 'react';

import Switch from '@mui/material/Switch';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Paper from '@mui/material/Paper';
import { useFetchUser, useGenerateTotpMfaSetup } from '../requests/userHooks';
import { SnackbarContext } from '../context/SnackbarContext';
import DisableMfaModal from '../components/DisableMfaModal';
import EnableMfaModal from '../components/EnableMfaModal';

function Account() {
  const { openErrorSnackbar }: any = useContext(SnackbarContext);
  const [checked, setChecked] = React.useState(false);
  const [totpImage, setTotpImage] = React.useState('');

  const { user, isLoadingUser } = useFetchUser();

  useEffect(() => {
    if (user?.isMfaEnabled) {
      setChecked(true);
    } else {
      setChecked(false);
    }
  }, [user]);

  const [openEnableMfaModal, setOpenEnableMfaModal] = React.useState(false);

  const handleCloseEnableMfaModal = () => {
    setOpenEnableMfaModal(false);
  };

  const handleClickOpenEnableMfaModal = () => {
    setOpenEnableMfaModal(true);
  };

  const { isLoadingGenerateTotpMfaSetup, generateTotpMfaSetup } = useGenerateTotpMfaSetup({
    onSuccess: (data) => {
      setTotpImage(data?.totpSetup?.image);
      handleClickOpenEnableMfaModal();
    },
    onError: (error) => {
      openErrorSnackbar(error?.response?.data?.message || error?.message);
    },
  });

  const [openDisableMfaModal, setOpenDisableMfaModal] = React.useState(false);

  const handleClickOpenDisableMfaModal = () => {
    setOpenDisableMfaModal(true);
  };

  const handleCloseDisableMfaModal = () => {
    setOpenDisableMfaModal(false);
  };
  const handleChangeMfaStatus = () => {
    if (!checked) {
      generateTotpMfaSetup();
    } else {
      handleClickOpenDisableMfaModal();
    }
  };

  return (
    <Box sx={{ minWidth: 120 }} className="mt-64 min-h-screen w-full flex flex-col items-center justify-center">
      <Paper className="p-32">
        <Typography variant="h5">{`Name: ${user?.name}`}</Typography>
        <Typography variant="h5" className="mt-12">{`Email:  ${user?.email}`}</Typography>
        <div className="flex items-center mt-12">
          <Typography variant="h5" className="mr-12">
            Mfa Status:
          </Typography>
          <Switch
            disabled={isLoadingGenerateTotpMfaSetup || isLoadingUser}
            checked={checked}
            onChange={handleChangeMfaStatus}
            inputProps={{ 'aria-label': 'controlled' }}
          />
        </div>
      </Paper>
      <DisableMfaModal handleClose={handleCloseDisableMfaModal} open={openDisableMfaModal} />
      <EnableMfaModal handleClose={handleCloseEnableMfaModal} open={openEnableMfaModal} totpImage={totpImage} />
    </Box>
  );
}

export default Account;
