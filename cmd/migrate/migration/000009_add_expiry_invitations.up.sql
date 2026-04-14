ALTER TABLE user_invitations 
ADD column expiry timestamp(0) with time zone not null;