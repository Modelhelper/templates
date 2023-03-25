select Id, StreamKey, Player = isnull(Player, 'jw') 
from LiveStreamCredentials
where roomid = (
        select roomid from BookingSession 
        where id = @id
    )    