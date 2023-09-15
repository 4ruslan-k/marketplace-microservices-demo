
import React from 'react';
import Card from '@mui/material/Card';
import Button from '@mui/material/Button';
import { bool, func, node, string } from 'prop-types';
import GitHub from '@mui/icons-material/GitHub';
import LoadingButton from '@mui/lab/LoadingButton';
import Image from 'next/image';
import config from '../config'
import googleIcon from '../assets/images/googleIcon.svg';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const githubIconStyles = {
  display: 'flex',
  justifyContent: 'flex-start',
  backgroundColor: 'rgba(51, 51, 51, 1)',
  color: 'white',
  '&:hover': {
    backgroundColor: 'rgba(51, 51, 51, 0.9)',
    boxShadow: 'none',
  },
};

function AuthCardWrapper({ children, handleSubmit, isValid, title, submitButtonLabel, isLoading }) {
  const pathname  = usePathname()

  return (
    <div className="w-full min-h-screen flex justify-center items-center mt-32">
      <Card className="w-400 p-32 flex flex-col justify-center">
        <h1 className="mb-32">{title}</h1>
        <form onSubmit={handleSubmit} className="flex flex-col justify-center items-center">
          <div className="w-full">{children}</div>
          <div className="mt-24 w-320 flex flex-col justify-center items-center">
            <LoadingButton
              className="w-250"
              variant="contained"
              color="primary"
              type="submit"
              disabled={!isValid || isLoading}
              loading={isLoading}
            >
              {submitButtonLabel}
            </LoadingButton>
            <div className="w-150 mt-16">
              <div className="flex flex-col">
                <div className="mb-4">
                  <Button
                    sx={{
                      root: {
                        display: 'flex',
                        justifyContent: 'flex-start',
                      },
                    }}
                    className="w-full flex justify-between"
                    href={`${config.API_URL}/v1/auth/social/google`}
                    startIcon={
                      <Image alt="google icon" style={{ height: 20, width: 20, marginRight: 'auto' }} src={googleIcon} />
                    }
                    variant="outlined"
                  >
                    Login with Google
                  </Button>
                </div>
                <div className="mb-4">
                  <Button
                    startIcon={<GitHub />}
                    sx={githubIconStyles}
                    className="w-full"
                    href={`${config.API_URL}/v1/auth/social/github`}
                    variant="contained"
                  >
                    Login with Github
                  </Button>
                </div>
                <div className="mb-4">
                {pathname === '/signin' ? (
                  <Link href="/signup" className="no-underline">
                    <Button
                      className="w-full"
                      variant="outlined"
                    >
                      Sign Up
                    </Button>
                  </Link>
                ):<Link href="/signin" className="no-underline">
                <Button
                  className="w-full"
                  variant="outlined"
                >
                  Sign In
                </Button>
              </Link>}
        
                </div>
              </div>
            </div>
          </div>
        </form>
      </Card>
    </div>
  );
}

AuthCardWrapper.propTypes = {
  children: node.isRequired,
  handleSubmit: func.isRequired,
  isValid: bool.isRequired,
  title: string.isRequired,
  submitButtonLabel: string.isRequired,
};

export default AuthCardWrapper;
