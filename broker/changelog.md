# 1.3.6
- Fixed problem in the lead function which would lose data sent from the frontend in case there was already a lead created.
- Now it also recovers the lead and updates it in the DB if the policy in the DB are in PartnershipLead or Lead status.
