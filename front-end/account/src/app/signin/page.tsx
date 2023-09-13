"use client";

import React, { useContext, useEffect, useState } from 'react';
import { useSearchParams,useRouter } from 'next/navigation';
import TextField from '@mui/material/TextField';
import * as yup from 'yup';
import { useFormik } from 'formik';

import { useFetchUser, useLoginUser } from '../requests/userHooks';
import { SnackbarContext } from '../context/SnackbarContext';
import AuthCardWrapper from '../components/AuthCardWrapper';

import MfaModal from '../components/MfaModal';
import { emailValidation, passwordValidation } from '../constants/validations';

function SignIn() {
  const { openSuccessSnackbar, openErrorSnackbar }: any = useContext(SnackbarContext);
  const [open, setOpen] = React.useState(false);
  const [passwordVerificationToken, setPasswordVerificationToken] = useState<string>('');

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const { refetchUser } = useFetchUser();
  const router = useRouter();
  const { isMfaFlow, token } = useSearchParams()

  useEffect(() => {
    const resetUrl = () => {
      router.push(router.pathname);
    };
    if (isMfaFlow && token) {
      setOpen(true);
      setPasswordVerificationToken(String(token));
      resetUrl();
    }
  }, [isMfaFlow, token, setOpen, setPasswordVerificationToken, router]);

  const { loginUser, isLoadingLogin } = useLoginUser({
    onSuccess: (data:any) => {
      if (data.user.isMfaEnabled) {
        handleClickOpen();
        setPasswordVerificationToken(data.user.passwordVerificationToken);
      } else {
        openSuccessSnackbar('Login successful');
        router.push('/');
        refetchUser();
      }
    },
    onError: (error:any) => {
      openErrorSnackbar(error?.response?.data?.message || error?.message);
    },
  });
  const signInSchema = yup.object().shape({
    email: emailValidation,
    password: passwordValidation,
  });

  const { errors, handleChange, touched, handleBlur, isValid, handleSubmit } = useFormik({
    initialValues: {
      email: '',
      password: '',
    },
    validationSchema: signInSchema,
    onSubmit: (values) => loginUser(values),
  });

  const emailError = errors.email && touched.email ? errors.email : null;
  const passwordError = errors.password && touched.password ? errors.password : null;

  return (
    <AuthCardWrapper
      isValid={isValid}
      handleSubmit={handleSubmit}
      title="Login To Your Account"
      submitButtonLabel="Login"
      isLoading={isLoadingLogin}
    >
      <div className="mb-10 w-full">
        <TextField
          className="w-full"
          label="Email"
          type="email"
          variant="outlined"
          name="email"
          id="email"
          onChange={handleChange}
          onBlur={handleBlur}
          helperText={emailError}
          error={Boolean(emailError)}
          required
        />
      </div>
      <TextField
        className="w-full"
        label="Password"
        type="password"
        variant="outlined"
        name="password"
        id="password"
        onChange={handleChange}
        onBlur={handleBlur}
        helperText={passwordError}
        error={Boolean(passwordError)}
        required
      />
      <MfaModal handleClose={handleClose} open={open} passwordVerificationToken={passwordVerificationToken} />
    </AuthCardWrapper>
  );
}

export default SignIn;
